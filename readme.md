# Go Wait Library with Generics

[![Build Status](https://github.com/ic-n/wait/workflows/continuous-integration/badge.svg)](https://github.com/ic-n/wait/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/ic-n/wait)](https://goreportcard.com/report/github.com/ic-n/wait)
[![Go Reference](https://pkg.go.dev/badge/github.com/ic-n/wait.svg)](https://pkg.go.dev/github.com/ic-n/wait)

The **Go Wait** library is a Go package that provides a simple and efficient way to manage goroutines and gather their results using generics. It is designed as a wrapper for the standard `sync.WaitGroup` and allows you to work with goroutines that may return values or errors. This library simplifies error handling and provides a convenient interface for gathering results from concurrent tasks.

## Installation

You can install the Go Wait library using the `go get` command:

```bash
go get github.com/ic-n/wait
```

## Getting Started

To use the Go Wait library in your Go project, import it into your code:

```go
import "github.com/ic-n/wait"
```

## Example

Here's an example of how to use the Go Wait library to run concurrent tasks and gather their results:

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ic-n/wait"
)

type Data struct {
	Headers struct {
		Host      string `json:"Host"`
		UserAgent string `json:"User-Agent"`
	} `json:"headers"`
	Origin string `json:"origin"`
	URL    string `json:"url"`
}

func main() {
	g := wait.New[Data]()

	g.Go(func(ctx context.Context) (Data, error) {
		rsp, err := http.Get("https://httpbin.org/get")
		if err != nil {
			return Data{}, err
		}

		var d Data
		if err := json.NewDecoder(rsp.Body).Decode(&d); err != nil {
			return Data{}, err
		}

		return d, nil
	})

	g.Go(func(ctx context.Context) (Data, error) {
		rsp, err := http.Post("https://httpbin.org/post", "text/plain", http.NoBody)
		if err != nil {
			return Data{}, err
		}

		var d Data
		if err := json.NewDecoder(rsp.Body).Decode(&d); err != nil {
			return Data{}, err
		}

		return d, nil
	})

	err := g.Gather(func(result Data) {
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
