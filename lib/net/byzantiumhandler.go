package net

import (
	"github.com/ssbc/common"
	"github.com/ssbc/lib/mysql"
	"sync"

	"github.com/cloudflare/cfssl/log"
	"time"
)

var (
	Nodes                 = 1 // 系统节点数
	Urls         []string = []string{"http://127.0.0.1:8000"}
	isSelfLeader bool     = true //leader
	blockState   BlockState
	voteCounts   int = 1
	Ports        string
	Sender       string = "windows"
	signatures   map[string][]byte
	senders      map[string]string
	transinblock int = 6000
	transtoredis int = 60000
	times        int = 0
	rounds       int = 10
	Testflag         = ""
)

type BlockState struct {
	sync.Mutex
	currBlock common.Block
	tmpBlock  common.Block
}

func (bs *BlockState) GetCurrB() common.Block {
	bs.Lock()
	defer bs.Unlock()
	return bs.currBlock
}

func (bs *BlockState) GetTmpB() common.Block {
	bs.Lock()
	defer bs.Unlock()
	return bs.tmpBlock
}

func (bs *BlockState) SetCurrB(b common.Block) {
	bs.Lock()
	defer bs.Unlock()
	bs.currBlock = b

}

func (bs *BlockState) SetTmpB(b common.Block) {
	bs.Lock()
	defer bs.Unlock()
	bs.tmpBlock = b
}

func (bs *BlockState) Checkblock(b *common.Block) bool {
	bs.Lock()
	defer bs.Unlock()
	if b.PrevHash != bs.currBlock.Hash {
		return false
	}
	if bs.tmpBlock.Hash == b.Hash {
		return false
	}
	return true

}

func (bs *BlockState) Checks(hash string) bool {
	bs.Lock()
	defer bs.Unlock()

	if bs.tmpBlock.Hash != hash {
		return true
	}
	if bs.tmpBlock.Hash == bs.currBlock.Hash {
		return true
	}
	return false

}

//func (bs *BlockState) StoreBlock() {
//	bs.Lock()
//	defer bs.Unlock()
//	log.Info("store the block into Mysql")
//	log.Info("Successfully stored the block", bs.tmpBlock)
//	common.Blockchains <- bs.tmpBlock
//	bs.currBlock = bs.tmpBlock
//
//}

func (bs *BlockState) CheckAndStore(hash string) {
	bs.Lock()
	defer bs.Unlock()

	if bs.tmpBlock.Hash != hash {
		log.Info("store_block: This round may finished. not equal to hash")
		return
	}
	if bs.tmpBlock.Hash == bs.currBlock.Hash {
		log.Info("store_block: This round may finished. equal to current")
		return
	}
	log.Info("Successfully stored the block, id: ", bs.tmpBlock.Id)
	//common.Blockchains <- bs.tmpBlock
	bs.currBlock = bs.tmpBlock

	blockId := mysql.InsertBlock(bs.currBlock) // 插入Block，并获取blockId
	bs.currBlock.Id = blockId
	mysql.InsertTransaction(bs.currBlock) // 插入Transaction

	t2 = time.Now()
	log.Info("耗时: ", t2.Sub(t1))
	//log.Info("times and len of blockchain: ", times+1, len(common.Blockchains))
	log.Info("---------------------------------------------------------------------------------------------------------------------------------------")
	if times+1 < rounds {
		times++
		//time.Sleep(time.Second)
		go SendTrans()
	}
}

func Init() {
	blockState.SetCurrB(common.B)
	signatures = make(map[string][]byte)
	senders = make(map[string]string)
	//log.Info("Byzantium Init Successfully")
}

func vote(s *Server) *serverEndpoint {
	return &serverEndpoint{
		Methods: []string{"POST"},
		Handler: voteHandler,
		Server:  s,
	}
}

func voteHandler(ctx *serverRequestContextImpl) (interface{}, error) {

	return nil, nil
}

func store_vote(v *Vote) {

}
