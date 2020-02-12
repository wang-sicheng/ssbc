package net

import (

	"github.com/ssbc/common"
	"sync"

	"github.com/cloudflare/cfssl/log"
)

var (
	Nodes = 2
	vm sync.Mutex
	voteCount int = 0
	revoCount int = 0
	isSelfLeader bool = true  //leader
	votes = sync.Map{}
	//votes = make(map[string]chan *Vote)
	revotes = sync.Map{}
	//revotes = make(map[string]chan *ReVote)
	//tmpBlock *common.Block
	//CurrentBlock common.Block
	blockState BlockState
	voteCounts int = 1
	Ports string
	Sender string = "windows"
)

type BlockState struct{
	sync.Mutex
	currBlock common.Block
	tmpBlock common.Block
}

func (bs *BlockState) GetCurrB()common.Block{
	bs.Lock()
	defer bs.Unlock()
	return bs.currBlock
}

func (bs *BlockState) GetTmpB()common.Block{
	bs.Lock()
	defer bs.Unlock()
	return bs.tmpBlock
}

func (bs *BlockState) SetCurrB(b common.Block){
	bs.Lock()
	defer bs.Unlock()
	bs.currBlock = b

}

func (bs *BlockState) SetTmpB(b common.Block){
	bs.Lock()
	defer bs.Unlock()
	bs.tmpBlock = b
}

func (bs *BlockState) Checkblock(b *common.Block)bool{
	bs.Lock()
	defer bs.Unlock()
	if b.PrevHash != bs.currBlock.Hash{
		return false
	}
	if bs.tmpBlock.Hash == b.Hash{
		return false
	}
	return true

}

func (bs *BlockState) Checks(hash string)bool{
	bs.Lock()
	defer bs.Unlock()

	if bs.tmpBlock.Hash != hash{
		return true
	}
	if bs.tmpBlock.Hash == bs.currBlock.Hash{
		return true
	}
	return false

}

func (bs *BlockState) StoreBlock(){
	bs.Lock()
	defer bs.Unlock()
	log.Info("store the block into Mysql")
	log.Info("Successfully stored the block", bs.tmpBlock)
	common.Blockchains <- bs.tmpBlock
	bs.currBlock = bs.tmpBlock

}

func Init(){
	blockState.SetCurrB(common.B)
	log.Info("Byzantium Init Successfully")
}








func vote(s *Server)*serverEndpoint{
	return &serverEndpoint{
		Methods: []string{ "POST"},
		Handler: voteHandler,
		Server:  s,
	}
}



func voteHandler(ctx *serverRequestContextImpl) (interface{}, error) {

	return nil, nil
}



func store_vote(v *Vote){

}









