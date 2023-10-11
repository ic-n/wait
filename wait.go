// Package wait allows you to work with goroutines that may return values or
// errors. This library simplifies error handling and provides a convenient
// interface for gathering results from concurrent tasks.
package wait

import (
	"context"
	"sync"
)

// Group wait group to process gorutines that return value and error.
type Group[T any] struct {
	ctx     context.Context
	cancel  context.CancelCauseFunc
	wg      sync.WaitGroup
	errOnce sync.Once
	results chan T
}

// New creates a new wait group with a background context.
func New[T any]() *Group[T] {
	return WithContext[T](context.Background())
}

// WithContext creates a new wait group with a parent context to control lifecycle.
func WithContext[T any](ctx context.Context) *Group[T] {
	gctx, cancel := context.WithCancelCause(ctx)

	return &Group[T]{
		ctx:    gctx,
		cancel: cancel,

		wg:      sync.WaitGroup{},
		results: make(chan T),
	}
}

// Go starts a new goroutine that executes the specified function.
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

// Gather gathers results from the executed goroutines. The function is
// called for each result, and any errors are collected. This function blocks
// until all goroutines have completed.
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

// Wait return results from the executed goroutines as slice and any errors
// encountered. This function blocks until all goroutines have completed.
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
