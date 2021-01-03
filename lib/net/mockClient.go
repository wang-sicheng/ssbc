package net

import (
	"encoding/json"
	"github.com/cloudflare/cfssl/log"
)

func mockClient(s *Server) *serverEndpoint {
	return &serverEndpoint{
		Methods:   []string{"GET", "POST", "HEAD"},
		Handler:   mockClientHandler,
		Server:    s,
		successRC: 200,
	}
}

// 模拟客户端向四个节点发送交易
func mockClientHandler(ctx *serverRequestContextImpl) (interface{}, error) {
	trans := generateTx()

	for _, tx := range trans {
		txJson, err := json.Marshal(tx)
		if err != nil {
			log.Info("transaction Marshal err: ", err)
		}
		Broadcast("receiveTx", txJson)
	}

	resp := TestInfoResponseNet{
		TName:   "hello",
		Version: "world",
	}
	return resp, nil
}
