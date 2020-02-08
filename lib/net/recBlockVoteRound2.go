package net

import (
	"github.com/cloudflare/cfssl/log"
	"encoding/json"
	"github.com/ssbc/lib/redis"
	rd "github.com/gomodule/redigo/redis"
	"github.com/ssbc/common"
	"time"
)

type ReVote struct{
	Sender string
	Vote []Vote
	Hash string
	V bool
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
	conn := redis.Pool.Get()
	defer conn.Close()
	_,err = conn.Do("SADD",v.Hash+"round2", b)
	if err != nil{
		log.Info("recBlockVoteRound2Handler err SADD: ", err)
	}
	vc,err := redis.ToInt(conn.Do("SCARD",v.Hash+"round2"))
	if err !=nil{
		log.Info("recBlockVoteRound1Hand2er err:", err)
	}
	log.Info("recBlockVoteRound1Hand2er revoteCount : ",vc)
	if vc == voteCounts{
		log.Info("statistic the votes")
		go statistic(v.Hash)
	}
	return nil, nil
}

func statistic(hash string){
	//statistic 2 round vote
	// then decide whether store the block or not

	conn := redis.Pool.Get()
	defer conn.Close()
	vs,err := rd.ByteSlices(conn.Do("SMEMBERS", hash+"round2"))
	if err != nil{
		log.Info("recBlockVoteRound2Handler err SMEMBERS: ", err)
	}
	revotes := []ReVote{}
	for _,data := range vs{
		t := ReVote{}
		err := json.Unmarshal(data, &t)
		if err != nil{
			log.Info("recBlockVoteRound2Handler err json: ", err)
			continue
		}
		revotes = append(revotes, t)
	}
	votecount := 0
	for _,data := range revotes{
		if data.V{
			votecount++
		}
	}
	if float64(votecount) > float64(Nodes)* 0.75{
		log.Info("recBlockVoteRound2Handler: vote round tow has received more than 3/4 affirmative votes")
		store_block(hash)
	}
}


func store_block(hash string){
	if tmpBlock.Hash == hash{
		log.Info("Pulling out tmpBlock")
	}
	log.Info("store the block")
	//lib.Db.insert(block)
	log.Info("store the block into Mysql")
	log.Info("Successfully stored the block")
	common.Blockchains <- *tmpBlock
	currentBlock = *tmpBlock
	t2 =time.Now()
	log.Info("duration: ",t2.Sub(t1))
	log.Info("times and len of blockchain: ", times+1, len(common.Blockchains))
	if times+2 != len(common.Blockchains){
		panic("mismatch")
	}
	if times < 199{
		times++
		SendTrans()
	}

}