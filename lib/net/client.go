package net

import (
	"encoding/json"
	"github.com/cloudflare/cfssl/log"
	"github.com/ssbc/common"
	"github.com/ssbc/crypto"
	"github.com/ssbc/lib/mysql"
	"strconv"
	"time"
)

type SendCoinsParams struct {
	From string `json:"from"`
	To string `json:"to"`
	Amount string `json:"amount"`
}

type GetTransParams struct {
	Address string `json:"address"`
	Limit string `json:"limit"`
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
	params := &SendCoinsParams{}
	err = json.Unmarshal(b, params)
	if err != nil {
		log.Info("ERR sendCoinsHandler: ", err)
	}
	amountInt, err := strconv.Atoi(params.Amount)
	if err != nil {
		return "amount type error, expect int", nil
	}
	//log.Infof("sendCoinsParams: %v", *params)
	senderInfo := mysql.QueryAccountInfo(params.From)
	//log.Infof("senderInfo from db: %v", senderInfo)
	if senderInfo != (common.Account{}){
		privateKey, publicKey := senderInfo.PrivateKey, senderInfo.PublicKey
		message := "coin amount "+params.Amount
		signature := crypto.SignECC([]byte(message), privateKey)
		newTrac := common.Transaction{
			SenderAddress:   params.From,
			ReceiverAddress: params.To,
			Timestamp:       time.Now().String(),
			Signature:       signature,
			SenderPublicKey: publicKey,
			TransferAmount:  amountInt,
			Message: message,
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

func getTransaction(s *Server) *serverEndpoint {
	return &serverEndpoint{
		Methods:   []string{"GET", "POST"},
		Handler:   getTransactionHandler,
		Server:    s,
		successRC: 200,
	}
}

// 交易查询
func getTransactionHandler(ctx *serverRequestContextImpl) (interface{}, error) {
	b, err := ctx.ReadBodyBytes()
	if err != nil {
		log.Info("ERR getTransactionHandler: ", err)
		return nil, err
	}
	log.Info("getTransactionHandler rec: ", string(b))
	params := &GetTransParams{}
	err = json.Unmarshal(b, params)
	if err != nil {
		log.Info("ERR getTransactionHandler: ", err)
		return nil, err
	}
	limitInt, err := strconv.Atoi(params.Limit)
	if err != nil {
		return "limit type error, expect int", nil
	}
	res := mysql.QueryTransInfo(params.Address, limitInt)
	return res, nil
}
