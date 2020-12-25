package net

import (
	"encoding/json"
	"github.com/cloudflare/cfssl/log"
	rd "github.com/gomodule/redigo/redis"
	"github.com/ssbc/common"
	"github.com/ssbc/lib/redis"
)

type ReVote struct {
	Sender string
	Vote   []Vote
	Hash   string
	V      bool
}

func recBlockVoteRound2(s *Server) *serverEndpoint {
	return &serverEndpoint{
		Methods: []string{"POST"},
		Handler: recBlockVoteRound2Handler,
		Server:  s,
	}
}

func recBlockVoteRound2Handler(ctx *serverRequestContextImpl) (interface{}, error) {
	b, err := ctx.ReadBodyBytes()
	if err != nil {
		log.Info("ERR recBlockVoteRound2Handler: ", err)
	}
	//log.Info("recBlockVoteRound2Handler: ", string(b))
	v := &ReVote{}
	err = json.Unmarshal(b, v)
	if err != nil {
		log.Info("ERR recBlockVoteRound2Handler: ", err)
	}
	conn := redis.Pool.Get()
	defer conn.Close()
	_, err = conn.Do("SADD", v.Hash+"round2", b)
	if err != nil {
		log.Info("recBlockVoteRound2Handler err SADD: ", err)
	}
	vc, err := redis.ToInt(conn.Do("SCARD", v.Hash+"round2"))
	if err != nil {
		log.Info("recBlockVoteRound2Handler err:", err)
	}
	log.Infof("收到第二轮投票 %d 张", vc)
	if vc == Nodes {
		go statistic(v.Hash)
	}
	return nil, nil
}

func statistic(hash string) {
	//statistic 2 round vote
	// then decide whether store the block or not
	if blockState.Checks(hash) {
		log.Info("store_block: This round may finished")
		return
	}
	conn := redis.Pool.Get()
	defer conn.Close()
	vs, err := rd.ByteSlices(conn.Do("SMEMBERS", hash+"round2"))
	if err != nil {
		log.Info("recBlockVoteRound2Handler err SMEMBERS: ", err)
	}
	revotes := []ReVote{}
	for _, data := range vs {
		t := ReVote{}
		err := json.Unmarshal(data, &t)
		if err != nil {
			log.Info("recBlockVoteRound2Handler err json: ", err)
			continue
		}
		revotes = append(revotes, t)
	}
	votecount := 0
	for _, data := range revotes {
		if data.V {
			votecount++
		}
	}
	//log.Info("同意票数: ", votecount)
	if votecount >= common.QuorumNumber(Nodes) {
		log.Infof("第二轮投票同意票数 %d 张，达到2f+1，准备存储区块", votecount)
		store_block(hash)
	} else {
		log.Infof("同意票数: %d, 不足2f+1，本轮结束，交易返还至Redis", votecount)
		restore_tx()
		log.Info("---------------------------------------------------------------------------------------------------------------------------------------")
	}
	if times+1 < rounds {
		times++
		//time.Sleep(time.Second)
		go SendTrans()
	}
}

func store_block(hash string) {
	blockState.CheckAndStore(hash)

}

func restore_tx() {
	blockState.restore_tx()
}
