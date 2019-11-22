package net

import (
	"fmt"
	"github.com/cloudflare/cfssl/log"
	"github.com/ssbc/common"
	"time"
)


var t2 = time.Now()
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
		t1 := time.Now()

		b :=[]byte(`{"BPM": 10}`)
		for i:=0 ; i< 100000 ; i++{
			Broadcast("recTransHash",b)
		}
		time.Sleep(time.Second)
		log.Info("blockchain len",len(common.Blockchain))
	log.Info("blockchain len",len(common.Blockchains))
	//log.Info("Now the newest 10 blocks is:")
	//l :=len(common.Blockchain)
	//log.Info("len of blockchain: " ,len(common.Blockchain))
	//for i:=0 ;l-1-i>=0&&i<10;i++{
	//	log.Info(common.Blockchain[i])
	//}
	time.Sleep(time.Second*5)
		log.Info("duration : ", t2.Sub(t1))



}