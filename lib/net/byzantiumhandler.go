package net

import (
	"encoding/json"
	"github.com/cloudflare/cfssl/log"
	"github.com/ssbc/common"
	"sync"
	"time"
)

var (
	vm sync.Mutex
	voteCount int = 0
	revoCount int = 0
	commonTrans int = 1
	votes = sync.Map{}
	//votes = make(map[string]chan *Vote)
	revotes = sync.Map{}
	//revotes = make(map[string]chan *ReVote)
	tmpBlock *common.Block
	voteCounts int = 1
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
	go findCommonTrans(message.BPM)
	return nil, nil
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

func recBlockVoteRound1(s *Server)*serverEndpoint{
	return &serverEndpoint{
		Methods: []string{ "POST"},
		Handler: recBlockVoteRound1Handler,
		Server:  s,
	}
}

func recBlockVoteRound1Handler(ctx *serverRequestContextImpl) (interface{}, error) {
	b,err := ctx.ReadBodyBytes()
	if err !=nil{
		log.Info("ERR recBlockVoteRound1Handler: ", err)
	}
	log.Info("recBlockVoteRound1Handler: ",string(b))
	v := &Vote{}
	err = json.Unmarshal(b, v)
	if err !=nil{
		log.Info("ERR recBlockVoteRound1Handler: ", err)
	}
	//if _,ok := votes[v.Hash];ok{
	var voteChan chan *Vote
	if vs,ok:=votes.Load(v.Hash);ok{
		voteChan,_ =vs.(chan *Vote)
		voteChan <- v
	}
	vc := len(voteChan)
		//votes[v.Hash] <- v
		// vc := len(votes[v.Hash])
	log.Info("recBlockVoteRound1Handler voteCount : ",vc)
	if vc == voteCounts{
		go voteForRoundNew(v.Hash)
	}
	//}

	//vm.Lock()
	//voteCount++
	//log.Info("recBlockVoteRound1Handler voteCount : ",voteCount)
	//if voteCount % 4 == 0{
	//	go voteForRound(v)
	//}
	//vm.Unlock()
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
	//revotes[v.Hash] <- v
	var revoteChan chan *ReVote
	if vs,ok:=revotes.Load(v.Hash);ok{
		revoteChan,_ =vs.(chan *ReVote)
		revoteChan <- v
	}
	vc := len(revoteChan)
	log.Info("recBlockVoteRound2Handler revoCount : ",vc)
	if vc == voteCounts{
		go statistic(v)
	}
	//vm.Lock()
	//revoCount++
	//log.Info("recBlockVoteRound2Handler revoCount: ",revoCount)
	//if revoCount % 4 == 0 {
	//	statistic(v)
	//}
	//vm.Unlock()
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
	Vote []*Vote
	Hash string
}

func voteHandler(ctx *serverRequestContextImpl) (interface{}, error) {

	return nil, nil
}

func voteForRound(v *Vote){
	//when receive whole nodes votes
	//then statistics


	store_vote(v)
	rv := &ReVote{Sender:"zhuanfa", Vote: []*Vote{v}}
	b, err := json.Marshal(rv)
	if err != nil{
		log.Info("voteForRound: ",err)
		return
	}

	Broadcast("recBlockVoteRound2",b)

}
func voteForRoundNew(hash string){
	//when receive whole nodes votes
	//then statistics

	vs:=[]*Vote{}
	var voteChan chan *Vote
	if vs,ok:=votes.Load(hash);ok{
		voteChan,_ =vs.(chan *Vote)

	}
	for i:=0;i<voteCounts;i++{
		vs = append(vs,<- voteChan)
	}
	//close(votes[hash])
	//delete(votes, hash)
	close(voteChan)
	votes.Delete(hash)
	rv := &ReVote{Sender:"zhuanfa", Vote: vs, Hash:hash}
	b, err := json.Marshal(rv)
	if err != nil{
		log.Info("voteForRound: ",err)
		return
	}
	revotes.Store(rv.Hash,make(chan *ReVote, 10))
	//revotes[rv.Hash] = make(chan *ReVote, 10)
	Broadcast("recBlockVoteRound2",b)

}

func store_vote(v *Vote){

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

	if commonTrans == 1 {
		Broadcast("recBlock",bb)
	}


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
		votes.Store(block.Hash,make(chan *Vote,10))
		//votes[block.Hash] = make(chan *Vote,10)
		tmpBlock = block
		v := &Vote{Sender : "hihihi",Hash : block.Hash, Vote : 1 }
		b, err := json.Marshal(v)
		if err != nil{
			log.Info("verify_block: ", err)
		}
		Broadcast("recBlockVoteRound1", b)
	}

}
func verify_block(block *common.Block)bool{
	return true
}

func statistic(rv *ReVote){
	//statistic 2 round vote
	// then decide whether store the block or not
	if tmpBlock.Hash == rv.Hash{
		log.Info("Pulling out tmpBlock")
	}
	var revoteChan chan *ReVote
	if vs,ok:=revotes.Load(rv.Hash);ok{
		revoteChan,_ =vs.(chan *ReVote)

	}
	//close(revotes[rv.Hash])
	//delete(revotes, rv.Hash)
	close(revoteChan)
	revotes.Delete(rv.Hash)
	log.Info("store the block")
	store_block()
	log.Info("Successfully stored the block")
	common.Blockchain = append(common.Blockchain, *tmpBlock)
	common.Blockchains <- *tmpBlock
	log.Info("Now the newest 10 blocks is:")
	l :=len(common.Blockchain)
	log.Info("len of blockchain: " ,len(common.Blockchain))
	for i := l-1 ; i >= 0 && l-1-i<10 ; i--{
		log.Info(common.Blockchain[i])
	}
	t2 =time.Now()







}


func store_block(){

	log.Info("store the block into Mysql")
	//lib.Db.insert(block)

}