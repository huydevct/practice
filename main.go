package main

import (
	"github.com/alitto/pond"
	"github.com/gin-gonic/gin"
	"net/http"
	call "practice/multiple-api"
)

func main() {
	urls := []string{
		"https://www.golang.org/",
		"https://www.google.com/",
		"https://www.github.com/",
	}
	reqs := []call.RequestApi{}
	for _, url := range urls {
		req := call.RequestApi{
			Url:    url,
			Method: "GET",
		}
		reqs = append(reqs, req)
	}

	pool := pond.New(1, 1000)
	defer pool.StopAndWait()

	res := make(chan []byte, len(reqs))
	resp := [][]byte{}
	r := gin.Default()
	r.GET("/mulFetch", func(c *gin.Context) {
		resp = call.MultipleCall(pool, reqs, res)

		c.JSON(http.StatusOK, gin.H{
			"response": resp,
		})
	})

	r.Run()
}
