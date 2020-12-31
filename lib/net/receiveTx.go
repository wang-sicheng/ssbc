package net

import (
	"encoding/json"
	"github.com/cloudflare/cfssl/log"
	"github.com/ssbc/common"
	"github.com/ssbc/crypto"
	"github.com/ssbc/lib/redis"
)

func receiveTx(s *Server) *serverEndpoint {
	return &serverEndpoint{
		Methods: []string{"POST"},
		Handler: receiveTxHandler,
		Server:  s,
	}
}

func receiveTxHandler(ctx *serverRequestContextImpl) (interface{}, error) {
	b, err := ctx.ReadBodyBytes()
	if err != nil {
		log.Info("ERR receiveTxHandler: ", err)
	}
	//log.Info("receiveBlockHandler rec: ", string(b))
	newTx := &common.Transaction{}
	err = json.Unmarshal(b, newTx)
	if err != nil {
		log.Info("ERR receiveBlockHandler newTx json: ", err)
	}
	//log.Info("receiveBlockHandler newTx: ", *newTx)
	if verifyTx(newTx) {
		CacheTx(b)
	}

	return nil, nil
}

func verifyTx(tran *common.Transaction) bool {
	res := crypto.VerifySignECC([]byte(tran.Message), tran.Signature, tran.SenderPublicKey)
	return res
}

//func CacheTx(newTx *common.Transaction){
//	if verifyTx(newTx){
//		//if docker.IsSmartContract(newTx){
//		//	log.Info("receiveBlockHandler: is SmartContract")
//		//	smi, err := docker.GenerateSCSpec(newTx)
//		//	if err != nil{
//		//		log.Info("ERR receiveTxHandler GenerateSCSpec: ", err)
//		//		return
//		//	}
//		//	smd,err := docker.Compile(smi)
//		//	if err != nil{
//		//		log.Info("ERR receiveTxHandler Compile: ", err)
//		//		return
//		//	}
//		//	log.Info("SmartComtractDefintion: ", *smd)
//		//	b,err := json.Marshal(smd)
//		//	if err != nil{
//		//		log.Info("ERR receiveTxHandler json smd: ", err)
//		//		return
//		//	}
//		//	newTx.Message = string(b)
//		//}
//
//		transbyte,err  := json.Marshal(newTx)
//		if err != nil{
//			log.Info("ERR receiveTxHandler json tx: ", err)
//			return
//		}
//		conn := redis.Pool.Get()
//		defer conn.Close()
//		_,err = conn.Do("RPUSH", "transPool", transbyte)
//		if err != nil{
//			log.Info("ERR receiveTxHandler RPUSH: ", err)
//		}
//	}
//}

func CacheTx(b []byte) {

	conn := redis.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("RPUSH", "transPool", b)
	if err != nil {
		log.Info("ERR receiveTxHandler RPUSH: ", err)
	}

	length, err2 := redis.ToInt(conn.Do("LLEN", "transPool"))
	if err2 != nil {
		log.Info("CacheTx LLEN err: ", err2)
	}
	//log.Infof("当前缓存池有交易：%d", length)
	if length >= 6000 {
		SendTrans()
	}

}
