package net

import (
	"github.com/ssbc/common"
	"github.com/ssbc/lib/mysql"


)

func newblockEndpoint(s *Server) *serverEndpoint {
	return &serverEndpoint{
		Methods: []string{"GET", "POST"},
		Handler: blockHandler,
		Server:  s,
	}
}

func blockHandler(ctx *serverRequestContextImpl) (interface{}, error) {
	var block common.Block
	err := ctx.ReadBody(&block)
	if err != nil{
		return nil,err
	}
	prevBlock := common.Blockchain[len(common.Blockchain)-1]
	newBlock := common.GenerateBlock(prevBlock, block)
	common.Blockchain = append(common.Blockchain, newBlock)





	mysql.InsertBlock(newBlock)
	resp := &blockResponseNet{}
	resp.Index = newBlock.Index
	resp.Hash = newBlock.Hash
	return resp, nil
}

type blockResponseNet struct {
	Index int

	Hash string
}

func newRecBlockRound2(s *Server)*serverEndpoint{
	return &serverEndpoint{
		Methods: []string{ "POST"},
		Handler: recBlockRound2Handler,
		Server:  s,
	}
}

func recBlockRound2Handler(ctx *serverRequestContextImpl) (interface{}, error) {

	return nil, nil
}

