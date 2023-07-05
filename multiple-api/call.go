package multiple_api

import (
	"context"
	"fmt"
	"github.com/alitto/pond"
	"net/http"
)

func fetch(url string, ctx context.Context, c chan<- *http.Response) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err == nil {
		resp.Body.Close()
	}
	c <- resp
	return err
}

func MultipleCall(pool *pond.WorkerPool, urls []string, c chan<- *http.Response) {
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

	//close(c)

	fmt.Println("Successfully fetched all URLs")
}
