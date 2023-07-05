package multiple_api

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type (
	Request struct {
		url string
	}

	QueryInterface interface {
		Query() Response
	}

	Response struct {
		Body string
		id   uint64
	}
)

// Query func: sample query url response a sample data
func (r *Request) Query() Response {
	id := rand.Uint64()
	body := r.url + strconv.FormatUint(id, 10)
	return Response{
		Body: body,
		id:   id,
	}
}

func NewRequest(url string) QueryInterface {
	return &Request{
		url: url,
	}
}

type (
	StreamRequest struct {
		Query QueryInterface
		Body  chan Response
		QuitC chan bool
	}

	StreamRequestInterface interface {
		query(allData chan Response)
		Quit(quit bool)
	}
)

func (sq *StreamRequest) Quit(quit bool) {
	sq.QuitC <- quit
}

func (sq *StreamRequest) query(allData chan Response) {
	startOneQuery := time.After(1 * time.Second)
	queryDone := make(chan bool, 1)
	timeOut := time.After(50 * time.Second)

	go func() {
		for {
			select {
			case <-startOneQuery:
				go func() {
					body := sq.Query.Query()
					sq.Body <- body
					queryDone <- true
				}()
			case <-queryDone:
				startOneQuery = time.After(1 * time.Second)
			case d := <-sq.Body:
				allData <- d
			case <-sq.QuitC:
				close(sq.Body)
				close(sq.QuitC)
				return
			case <-timeOut:
				close(sq.QuitC)
				close(sq.Body)
				return
			}
		}
	}()
}

func NewStreamRequest(Query QueryInterface) StreamRequestInterface {
	return &StreamRequest{
		Query: Query,
		Body:  make(chan Response),
		QuitC: make(chan bool),
	}
}

type (
	MergerAllStreamData struct {
		AllStream     []StreamRequestInterface
		AllData       chan Response
		QuitAllStream chan bool
	}
	MergerAllStreamDataInterface interface {
		Merge()
	}
)

func (m MergerAllStreamData) Merge() {
	// start all stream
	for _, v := range m.AllStream {
		v.query(m.AllData)
	}
	timeOut := time.After(10 * time.Second)

	//merge data
	go func() {
		for {
			select {
			case d := <-m.AllData:
				fmt.Printf("all Data Stream, %v \n", d)
			case <-m.QuitAllStream:
				fmt.Println("send Quit ALl channel")
				for _, v := range m.AllStream {
					v.Quit(true)
				}
				// todo case send quit one channel
			case <-timeOut:
				fmt.Println("send Quit ALl channel")
				for _, v := range m.AllStream {
					v.Quit(true)
				}
			}
		}
	}()
}

func NewMergerAllStreamData(listStream ...StreamRequestInterface) MergerAllStreamDataInterface {
	var allStream []StreamRequestInterface
	for _, one := range listStream {
		allStream = append(allStream, one)
	}

	return &MergerAllStreamData{
		AllStream:     allStream,
		AllData:       make(chan Response),
		QuitAllStream: make(chan bool),
	}
}
