package net

import (
	"fmt"
	"time"
)

type TestInfo struct {

	TName string

	Version string
}

func newTestInfoEndpoint(s *Server) *serverEndpoint {
	return &serverEndpoint{
		Methods: []string{"GET", "POST", "HEAD"},
		Handler: testinfoHandler,
		Server:  s,
	}
}

type TestInfoResponseNet struct {
	TName string

	Version string
}
func testinfoHandler(ctx *serverRequestContextImpl) (interface{}, error) {
	b,err := ctx.ReadBodyBytes()
	go SendTrans()
	fmt.Println(string(b),err)
	resp := TestInfoResponseNet{
		TName: "hello",
		Version: "world",
	}
	return resp, nil
}

func SendTrans(){

		b :=[]byte(`{"BPM": 10}`)
		Broadcast("recTransHash",b)
		time.Sleep(time.Second*5)


}