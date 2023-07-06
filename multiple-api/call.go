package multiple_api

import (
	"context"
	"fmt"
	"github.com/alitto/pond"
	"io"
	"net/http"
)

func fetch(url string, ctx context.Context, c chan []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err == nil {
		defer resp.Body.Close()
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	c <- body
	return err
}

func MultipleCall(pool *pond.WorkerPool, urls []string, c chan []byte) [][]byte {
	group, ctx := pool.GroupContext(context.Background())

	for _, url := range urls {
		group.Submit(func() error {
			fmt.Println("start fetch")
			err := fetch(url, ctx, c)
			return err
		})
	}

	err := group.Wait()
	if err != nil {
		fmt.Printf("Failed to fetch URLs: %v", err)
	}

	res := [][]byte{}

	res = append(res, <-c)

	fmt.Println("Successfully fetched all URLs")

	return res
}
