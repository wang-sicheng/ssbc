package net

import (
	"github.com/cloudflare/cfssl/log"
	"github.com/ssbc/common"
	"encoding/json"
	"github.com/ssbc/lib/redis"
	rd "github.com/gomodule/redigo/redis"
)

func receiveBlock(s *Server)*serverEndpoint{
	return &serverEndpoint{
		Methods: []string{ "POST"},
		Handler: receiveBlockHandler,
		Server:  s,
	}
}

func receiveBlockHandler(ctx *serverRequestContextImpl) (interface{}, error) {
	b,err := ctx.ReadBodyBytes()
	if err !=nil{
		log.Info("ERR receiveBlockHandler: ", err)
	}
	log.Info("receiveBlockHandler: ",string(b))
	newBlock := &common.Block{}
	err = json.Unmarshal(b, newBlock)
	if err !=nil{
		log.Info("ERR receiveBlockHandler: ", err)
	}
	log.Info("receiveBlockHandler: ", newBlock)
	if !blockState.Checkblock(newBlock){
		log.Info("receiveBlockHandler: Hash mismatch. This round may finish")
		return nil, nil
	}
	go verify(newBlock)
	return nil, nil
}

func verify(block *common.Block){
	// sender
	//tmpblock change to cache to redis next version
	voteBool := false
	if verify_block(block){
		voteBool = true
	}
	log.Info("verify block: ",voteBool)
	blockState.SetTmpB(*block)
	v := &Vote{Sender : Sender, Hash : block.Hash, Vote : voteBool }
	b, err := json.Marshal(v)
	if err != nil{
		log.Info("verify_block: ", err)
	}
	log.Info("vote: ", string(b))
	Broadcast("recBlockVoteRound1", b)
}


func verify_block(block *common.Block)bool{
	//验证逻辑 验签 验证交易 验证merkle tree root
	currentBlock := blockState.GetCurrB()
	if block.PrevHash != currentBlock.Hash{
		log.Info("This round may finish")
		return false
	}
	if block.Signature != "Signature"{
		log.Info("verify block: Signature mismatch")
		return false
	}
	return verifyBlockTx(block,&currentBlock)
}

func verifyBlockTx(b *common.Block, currentBlock *common.Block)bool{

	transCache := []interface{}{"verifyBlockTxCache"+b.Hash}
	for _,data := range b.TX{
		b,err := json.Marshal(data)
		if err != nil{
			log.Info("verifyBlockTx json err: ", err)
		}
		transCache = append(transCache, b)
	}
	if len(transCache) == 1{
		transCache = append(transCache, []byte{})
	}
	conn := redis.Pool.Get()
	defer conn.Close()
	_,err := conn.Do("SADD", transCache...)
	if err != nil{
		log.Info("verifyBlockTx err SADD: ", err)
	}
	commonTrans,err := rd.Strings(conn.Do("SINTER", "verifyBlockTxCache"+b.Hash, "CommonTxCache4verify"+ currentBlock.Hash))
	if err !=nil{
		log.Info("verifyBlockTx err SINTER: ", err)
	}
	log.Info("verifyBlockTx commonTrans:   ", commonTrans)
	log.Info("verifyBlockTx len trans commonTrans :   ", len(b.TX), len(commonTrans))
	if len(b.TX) != len(commonTrans){
		return false
	}
	return true
}