# Go Wait Library with Generics

The **Go Wait** library is a Go package that provides a simple and efficient way to manage goroutines and gather their results using generics. It is designed as a wrapper for the standard `sync.WaitGroup` and allows you to work with goroutines that may return values or errors. This library simplifies error handling and provides a convenient interface for gathering results from concurrent tasks.

## Installation

You can install the Go Wait library using the `go get` command:

```bash
go get github.com/your/repository/wait
```

## Getting Started

To use the Go Wait library in your Go project, import it into your code:

```go
import "github.com/your/repository/wait"
```

## Usage

The Go Wait library provides a `Group[T]` data structure with the following methods:

### `New[T any]() *Group[T]`

Creates a new `Group` with a background context.

```go
group := wait.New[int]()
```

### `WithContext[T any](ctx context.Context) *Group[T]`

Creates a new `Group` with a custom context. You can use this to control the lifecycle of the group.

```go
ctx := context.Background()
group := wait.WithContext[int](ctx)
```

### `Go(fn func(ctx context.Context) (T, error))`

Starts a new goroutine that executes the specified function. The function takes a context and returns a value of type `T` and an error.

```go
group.Go(func(ctx context.Context) (int, error) {
    // Your concurrent task logic here
    return 42, nil
})
```

### `Gather(gatherer func(T)) error`

Gathers results from the executed goroutines. The `gatherer` function is called for each result, and any errors are collected. This function blocks until all goroutines have completed.

```go
err := group.Gather(func(result int) {
    // Process the result
})
```

### `Wait() ([]T, error)`

Gathers results from the executed goroutines and returns a slice of results and any errors encountered. This function also blocks until all goroutines have completed.

```go
results, err := group.Wait()
```

## Example

Here's an example of how to use the Go Wait library to run concurrent tasks and gather their results:

```go
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
```

In this example, the Go Wait library is used to execute concurrent tasks and gather their results. It simplifies the handling of concurrent tasks and error management.

## License

This library is open-source and released under the [MIT License](LICENSE).
