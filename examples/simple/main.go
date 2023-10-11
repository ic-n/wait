package main

import (
	"context"
	"fmt"

	"github.com/ic-n/wait"
)

func main() {
	group := wait.New[int]()

	for i := 0; i < 5; i++ {
		i := i
		group.Go(func(ctx context.Context) (int, error) {
			return i * 2, nil
		})
	}

	err := group.Gather(func(result int) {
		fmt.Println("Result:", result)
	})

	if err != nil {
		fmt.Println("Error:", err)
	}
}
