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

func verifyTx(tx *common.Transaction) bool {
	res := crypto.VerifySignECC([]byte(tx.Message), tx.Signature, tx.SenderPublicKey)
	return res
}

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
	if length >= 6000 && !Processing {
		SendTrans()
		Processing = true
	}

}
