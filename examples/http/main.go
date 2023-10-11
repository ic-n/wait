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
	group := wait.New[Data]()

	group.Go(func(ctx context.Context) (Data, error) {
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

	group.Go(func(ctx context.Context) (Data, error) {
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

	err := group.Gather(func(result Data) {
		fmt.Println("Result:", result)
	})

	if err != nil {
		fmt.Println("Error:", err)
	}
}
