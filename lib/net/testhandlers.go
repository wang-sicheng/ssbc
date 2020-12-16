package net

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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

// 区块链测试接口
func testinfoHandler(ctx *serverRequestContextImpl) (interface{}, error) {

	log.Info("ctx.req.RemoteAddr: ", ctx.req.RemoteAddr)
	b, err := ctx.ReadBodyBytes()
	if err != nil {
		log.Info("ERR receiveTxHandler: ", err)
	}
	log.Info("receiveBlockHandler 接收的请求内容：\n", string(b))
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

// 广播TransHash
func SendTrans() {

	if flag {
		//flushall()
		//time.Sleep(time.Second)
		recTrans()
		t1 = time.Now()
		flag = false
	}
	a := pullTrans()
	transhash := TransHash{}
	transhash.BlockHash = blockState.GetCurrB().Hash
	m := make(map[string][]byte)
	transCache4verify := []interface{}{"CommonTxCache4verify" + transhash.BlockHash}
	for _, data := range a {
		hash := sha256.Sum256(data)
		hashString := hex.EncodeToString(hash[:])
		transhash.TransHashs = append(transhash.TransHashs, hashString)
		m[hashString] = data
		transCache4verify = append(transCache4verify, data)
	}
	b, err := json.Marshal(transhash)
	if err != nil {
		log.Info("test err: ", err)
		return
	}
	mb, err := json.Marshal(m)
	if err != nil {
		log.Info("test err m: ", err)
		return
	}

	conn := redis.Pool.Get()
	defer conn.Close()
	_, err = conn.Do("SET", "CommonTxCache"+transhash.BlockHash, mb)
	if err != nil {
		log.Info("test err SET: ", err)
	}
	_, err = conn.Do("SADD", transCache4verify...)
	if err != nil {
		log.Info("test err SADD: ", err)
	}
	//log.Info("common.Blockchains长度: ", len(common.Blockchains))
	log.Infof("广播 %d 条交易的Hash", len(transhash.TransHashs))
	Broadcast("recTransHash", b)

}

func Flushall() {
	conn := redis.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("flushall")
	if err != nil {
		panic(err)
	}
	log.Info("Redis已清空")

}
