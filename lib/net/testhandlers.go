package net

import (
	"github.com/cloudflare/cfssl/log"
	"time"
	"encoding/hex"
	"encoding/json"
	"github.com/ssbc/common"
	"crypto/sha256"
	"github.com/ssbc/lib/redis"
)

var(
	t1 time.Time
	t2 time.Time
	flag bool = true
	times int = 0
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
		successRC: 200,
	}
}

type TestInfoResponseNet struct {
	TName string

	Version string
}
func testinfoHandler(ctx *serverRequestContextImpl) (interface{}, error) {


	log.Info("ctx.req.RemoteAddr: ",ctx.req.RemoteAddr)



	go SendTrans()

	resp := TestInfoResponseNet{
		TName: "hello",
		Version: "world",
	}
	return resp, nil
}

func SendTrans(){

		//b,err := json.Marshal(common.Tx100)
		//if err != nil{
		//	log.Info("test err : ",err)
		//}
		//for i:=0 ; i< 1 ; i++{
		//	Broadcast("recTransHash",b)
		//}
		//time.Sleep(time.Second)
	//trans := make(chan []byte,100)
	//
	//go transToRedis(trans)
	//for i:=0;i<100;i++{
	//	marshalTrans(trans)
	//}
	//for i:=0;i<10;i++{
	//	recTrans()
	//}
	//log.Info("bye")
	//return
	recTrans()
	if flag{

		t1 = time.Now()
		flag = false
	}

	a := pullTrans()
	transhash := TransHash{}
	transhash.BlockHash = currentBlock.Hash
	m := make(map[string][]byte)
	transCache4verify := []interface{}{"CommonTxCache4verify"}
	for _,data := range a{
		hash := sha256.Sum256(data)
		hashString := hex.EncodeToString(hash[:])
		transhash.TransHashs = append(transhash.TransHashs, hashString)
		m[hashString] = data
		transCache4verify = append(transCache4verify, data)
	}
	b,err := json.Marshal(transhash)
	if err !=nil{
		log.Info("test err: ", err)
		return
	}
	mb,err := json.Marshal(m)
	if err !=nil{
		log.Info("test err m: ", err)
		return
	}

	conn := redis.Pool.Get()
	defer conn.Close()
	_,err = conn.Do("SET", "CommonTxCache"+currentBlock.Hash, mb)
	if err != nil{
		log.Info("test err SET: ", err)
	}
	_,err = conn.Do("SADD", transCache4verify...)
	if err != nil{
		log.Info("test err SADD: ", err)
	}

	Broadcast("recTransHash",b)
	log.Info("blockchain len",len(common.Blockchains))


}