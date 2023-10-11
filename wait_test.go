package wait_test

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"testing"
	"time"

	"github.com/ic-n/wait"
	"golang.org/x/exp/slices"
)

type (
	routine struct {
		ms  int
		v   string
		err error
	}
	testCase struct {
		routines []routine
		exp      []string
		expErr   error
	}
)

var testGroupCases = []testCase{
	{
		routines: []routine{
			{ms: 100, v: "1"},
			{ms: 200, v: "2"},
		},
		exp: []string{"1", "2"},
	},
	{
		routines: []routine{
			{ms: 100, v: "1"},
			{ms: 100, v: "2"},
			{ms: 100, v: "3"},
		},
		exp: []string{"1", "2", "3"},
	},
	{
		routines: []routine{
			{ms: 100, v: "1"},
			{ms: 200, err: fs.ErrNotExist},
			{ms: 300, v: "3"},
		},
		expErr: fs.ErrNotExist,
	},
}

func TestGroup(t *testing.T) {
	t.Parallel()

	for i, tc := range testGroupCases {
		t.Run(fmt.Sprintf("test-case-wait-%d", i), func(t *testing.T) {
			g := wait.New[string]()

			for _, r := range tc.routines {
				r := r
				g.Go(func(ctx context.Context) (string, error) {
					time.Sleep(time.Duration(r.ms) * time.Millisecond)
					return r.v, r.err
				})
			}

			result, err := g.Wait()
			slices.Sort(result) // for tests consistency

			if tc.exp != nil && !slices.Equal(tc.exp, result) {
				t.Fatalf("expected %s, got %v", tc.exp, result)
			}
			if !errors.Is(tc.expErr, err) {
				t.Fatalf("expected %v, got %v", tc.expErr, err)
			}
		})
		t.Run(fmt.Sprintf("test-case-gather-%d", i), func(t *testing.T) {
			g := wait.New[string]()

			for _, r := range tc.routines {
				r := r
				g.Go(func(ctx context.Context) (string, error) {
					time.Sleep(time.Duration(r.ms) * time.Millisecond)
					return r.v, r.err
				})
			}

			result := []string{}
			err := g.Gather(func(s string) {
				result = append(result, s)
			})
			slices.Sort(result) // for tests consistency

			if tc.exp != nil && !slices.Equal(tc.exp, result) {
				t.Fatalf("expected %s, got %v", tc.exp, result)
			}
			if !errors.Is(tc.expErr, err) {
				t.Fatalf("expected %v, got %v", tc.expErr, err)
			}
		})
	}
}
