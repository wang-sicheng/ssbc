package net

import (
	"encoding/json"
	"github.com/cloudflare/cfssl/log"
	"github.com/ssbc/common"
	"github.com/ssbc/lib/mysql"
	"github.com/ssbc/crypto"
	"time"
	"unsafe"
)

type sendCoinsParams struct {
	from string
	to string
	amount int
}

func sendCoins(s *Server) *serverEndpoint {
	return &serverEndpoint{
		Methods:   []string{"GET", "POST"},
		Handler:   sendCoinsHandler,
		Server:    s,
		successRC: 200,
	}
}

// 模拟client交易签名，构建交易
func sendCoinsHandler(ctx *serverRequestContextImpl) (interface{}, error) {
	b, err := ctx.ReadBodyBytes()
	if err != nil {
		log.Info("ERR sendCoinsHandler: ", err)
	}
	log.Info("sendCoinsHandler rec: ", string(b))
	params := &sendCoinsParams{}
	err = json.Unmarshal(b, params)
	if err != nil {
		log.Info("ERR sendCoinsHandler: ", err)
	}
	senderInfo := mysql.QueryAccountInfo(params.from)
	if senderInfo != (common.Account{}){
		privateKey, publicKey := senderInfo.PrivateKey, senderInfo.PublicKey
		signature := crypto.SignECC(Int2Byte(params.amount), privateKey)
		newTrac := common.Transaction{
			SenderAddress:   params.from,
			ReceiverAddress: params.to,
			Timestamp:       time.Now().String(),
			Signature:       signature,
			SenderPublicKey: publicKey,
			TransferAmount:  params.amount,
		}
		b, err := json.Marshal(&newTrac)
		if err != nil {
			log.Error("json marshal error")
			return nil, err
		}
		Broadcast("receiveTx", b)
	}else {
		return "address error", nil
	}
	return nil, nil
}

func Int2Byte(data int)(ret []byte){
	var l = unsafe.Sizeof(data)
	ret = make([]byte, l)
	var tmp = 0xff
	var index uint = 0
	for index=0; index<uint(l); index++{
		ret[index] = byte((tmp<<(index*8) & data)>>(index*8))
	}
	return ret
}

func Byte2Int(data []byte)int{
	var ret = 0
	var l = len(data)
	var i uint = 0
	for i=0; i<uint(l); i++{
		ret = ret | (int(data[i]) << (i*8))
	}
	return ret
}
