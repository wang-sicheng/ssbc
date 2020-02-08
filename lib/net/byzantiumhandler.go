package net

import (

	"github.com/ssbc/common"
	"sync"

)

var (
	Nodes = 1
	vm sync.Mutex
	voteCount int = 0
	revoCount int = 0
	isSelfLeader bool = true  //leader
	votes = sync.Map{}
	//votes = make(map[string]chan *Vote)
	revotes = sync.Map{}
	//revotes = make(map[string]chan *ReVote)
	tmpBlock *common.Block
	currentBlock common.Block
	voteCounts int = 1
	Ports string
)











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









