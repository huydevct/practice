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

	pool := pond.New(1, 1000)
	defer pool.StopAndWait()

	res := make(chan *http.Response, len(urls))

	r := gin.Default()
	r.GET("/mulFetch", func(c *gin.Context) {
		go call.MultipleCall(pool, urls, res)

		c.JSON(http.StatusOK, gin.H{
			"response": "pong",
		})
	})

	r.Run()
}
