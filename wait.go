package wait

import (
	"context"
	"sync"
)

type Group[T any] struct {
	ctx     context.Context
	cancel  context.CancelCauseFunc
	wg      sync.WaitGroup
	errOnce sync.Once
	results chan T
}

func New[T any]() *Group[T] {
	return WithContext[T](context.Background())
}

func WithContext[T any](ctx context.Context) *Group[T] {
	gctx, cancel := context.WithCancelCause(ctx)

	return &Group[T]{
		ctx:    gctx,
		cancel: cancel,

		wg:      sync.WaitGroup{},
		results: make(chan T),
	}
}

func (g *Group[T]) Go(fn func(ctx context.Context) (T, error)) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()

		v, err := fn(g.ctx)
		if err != nil {
			g.errOnce.Do(func() {
				g.cancel(err)
			})
			return
		}

		g.results <- v
	}()
}

func (g *Group[T]) Gather(gatherer func(T)) error {
	defer g.errOnce.Do(func() {
		g.cancel(nil)
	})

	readDone := make(chan struct{}, 1)

	go func() {
		defer func() { readDone <- struct{}{} }()
		for {
			select {
			case v, ok := <-g.results:
				if !ok {
					return
				}
				gatherer(v)
			case <-g.ctx.Done():
				return
			}
		}
	}()

	go func() {
		g.wg.Wait()
		close(g.results)
	}()

	<-readDone

	err := context.Cause(g.ctx)
	if err != nil {
		return err
	}

	return nil
}

func (g *Group[T]) Wait() ([]T, error) {
	defer g.errOnce.Do(func() {
		g.cancel(nil)
	})

	vs := make([]T, 0)
	readDone := make(chan struct{}, 1)

	go func() {
		defer func() { readDone <- struct{}{} }()
		for {
			select {
			case v, ok := <-g.results:
				if !ok {
					return
				}
				vs = append(vs, v)
			case <-g.ctx.Done():
				return
			}
		}
	}()

	go func() {
		g.wg.Wait()
		close(g.results)
	}()

	<-readDone

	err := context.Cause(g.ctx)
	if err != nil {
		return nil, err
	}

	return vs, nil
}
