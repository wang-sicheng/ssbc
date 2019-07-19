package net

import (
	"encoding/json"
	"github.com/cloudflare/cfssl/log"
	"github.com/ssbc/common"
)

func receive_trans_bitarry(s *Server)*serverEndpoint{
	return &serverEndpoint{
		Methods: []string{ "POST"},
		Handler: rtbitarryHandler,
		Server:  s,
	}
}

func rtbitarryHandler(ctx *serverRequestContextImpl) (interface{}, error) {
	b,err := ctx.ReadBodyBytes()
	message := &common.Message{}
	err = json.Unmarshal(b, message)
	if err != nil{
		log.Info(err)
	}
	log.Info("rtbitarryHandler: ", string(b))
	findCommonTrans(message.BPM)
	return nil, nil
}

func recBlockVoteRound1(s *Server)*serverEndpoint{
	return &serverEndpoint{
		Methods: []string{ "POST"},
		Handler: recBlockVoteRound1Handler,
		Server:  s,
	}
}

func recBlockVoteRound1Handler(ctx *serverRequestContextImpl) (interface{}, error) {

	return nil, nil
}

func recBlockVoteRound2(s *Server)*serverEndpoint{
	return &serverEndpoint{
		Methods: []string{ "POST"},
		Handler: recBlockVoteRound2Handler,
		Server:  s,
	}
}

func recBlockVoteRound2Handler(ctx *serverRequestContextImpl) (interface{}, error) {
	b,err := ctx.ReadBodyBytes()
	if err !=nil{
		log.Info("ERR recBlockVoteRound2Handler: ", err)
	}
	log.Info("recBlockVoteRound2Handler: ",string(b))
	v := &ReVote{}
	err = json.Unmarshal(b, v)
	if err !=nil{
		log.Info("ERR recBlockVoteRound2Handler: ", err)
	}
	log.Info("recBlockVoteRound2Handler: ",v)

	return nil, nil
}

func vote(s *Server)*serverEndpoint{
	return &serverEndpoint{
		Methods: []string{ "POST"},
		Handler: voteHandler,
		Server:  s,
	}
}

type ReVote struct{
	Sender string
	Vote *Vote
}

func voteHandler(ctx *serverRequestContextImpl) (interface{}, error) {
	b,err := ctx.ReadBodyBytes()
	if err !=nil{
		log.Info("ERR voteHandler: ", err)
	}
	log.Info("voteHandler: ",string(b))
	v := &Vote{}
	err = json.Unmarshal(b, v)
	if err !=nil{
		log.Info("ERR voteHandler: ", err)
	}
	log.Info("voteHandler: ",v)
	go voteForRound(v)
	return nil, nil
}

func voteForRound(v *Vote){
	//when receive whole nodes votes
	//then statistics


	store_vote(v)
	rv := &ReVote{Sender:"zhuanfa", Vote: v}
	b, err := json.Marshal(rv)
	if err != nil{
		log.Info("recBlockVoteRound: ",err)
		return
	}
	Broadcast("recBlockVoteRound",b)




}

func store_vote(v *Vote){

}
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
	go verify(newBlock)
	return nil, nil
}

func findCommonTrans(bpm int){
	//do something
	b := common.Block{BPM:bpm}
	b = common.GenerateBlock(common.Blockchain[len(common.Blockchain)-1],b)
	bb, err := json.Marshal(b)
	if err != nil{
		log.Info("ERR findCommonTrans: ", err)
		return
	}
	Broadcast("recBlock",bb)

}
type Vote struct{
	Sender string
	Hash string
	Vote int
}
func verify(block *common.Block){
	// do something like verify block
	// now we consider verify
	if verify_block(block){
		v := &Vote{Sender : "hihihi",Hash : block.Hash, Vote : 1 }
		b, err := json.Marshal(v)
		if err != nil{
			log.Info("verify_block: ", err)
		}
		Broadcast("recVote1", b)
	}

}
func verify_block(block *common.Block)bool{
	return true
}