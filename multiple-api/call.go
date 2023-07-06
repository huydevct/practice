package multiple_api

import (
	"bytes"
	"context"
	"fmt"
	"github.com/alitto/pond"
	"io"
	"net/http"
)

type RequestApi struct {
	Url    string
	Method string
	Body   string
}

func fetch(req RequestApi, c chan []byte) error {
	reqBody := bytes.NewReader([]byte(req.Body))
	request, err := http.NewRequest(req.Method, req.Url, reqBody)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return err
	}

	resp, err := http.DefaultClient.Do(request)
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

func MultipleCall(pool *pond.WorkerPool, reqs []RequestApi, c chan []byte) [][]byte {
	group, _ := pool.GroupContext(context.Background())

	for _, req := range reqs {
		group.Submit(func() error {
			fmt.Println("start fetch")
			err := fetch(req, c)
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
