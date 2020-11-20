package net

import (
	"github.com/cloudflare/cfssl/log"
	"github.com/ssbc/lib/redis"
	"time"
)

var (
	t1   time.Time
	t2   time.Time
	flag bool = true
)

type TestInfo struct {
	TName string

	Version string
}

func newTestInfoEndpoint(s *Server) *serverEndpoint {
	return &serverEndpoint{
		Methods:   []string{"GET", "POST", "HEAD"},
		Handler:   testinfoHandler,
		Server:    s,
		successRC: 200,
	}
}

type TestInfoResponseNet struct {
	TName string

	Version string
}

func testinfoHandler(ctx *serverRequestContextImpl) (interface{}, error) {

	log.Info("ctx.req.RemoteAddr: ", ctx.req.RemoteAddr)
	b, err := ctx.ReadBodyBytes()
	if err != nil {
		log.Info("ERR receiveTxHandler: ", err)
	}
	log.Info("receiveBlockHandler: ", string(b))
	//	newTx := &common.Transaction{}
	//	err = json.Unmarshal(b, newTx)
	//	if err !=nil{
	//		log.Info("ERR receiveBlockHandler newTx json: ", err)
	//	}
	//	s := `package main
	//
	//import (
	//	"fmt"
	//
	//)
	//
	//func main() {
	//	fmt.Println("Hello World")
	//}
	//	`
	//	smi := &docker.SmartContractInit{"TTEESSTT", "windows", "1.0", []byte(s)}
	//	b,_ = json.Marshal(smi)
	//	newTx.Message = string(b)
	//	b,_ = json.Marshal(newTx)
	//	go Broadcast("receiveTx", b)
	//	return nil,nil
	go SendTrans()
	resp := TestInfoResponseNet{
		TName:   "hello",
		Version: "world",
	}
	return resp, nil
}

func SendTrans() {

	if flag {
		//flushall()
		//time.Sleep(time.Second)

		t1 = time.Now()
		flag = false
	}

}

func Flushall() {
	conn := redis.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("flushall")
	if err != nil {
		panic(err)
	}
	log.Info("flushall success")

}
