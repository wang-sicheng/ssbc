package net

import (
	"encoding/json"
	"github.com/cloudflare/cfssl/log"
	rd "github.com/gomodule/redigo/redis"
	"github.com/ssbc/lib/redis"
)

type Vote struct {
	Sender string
	Hash   string
	Vote   bool
}

func recBlockVoteRound1(s *Server) *serverEndpoint {
	return &serverEndpoint{
		Methods: []string{"POST"},
		Handler: recBlockVoteRound1Handler,
		Server:  s,
	}
}

func recBlockVoteRound1Handler(ctx *serverRequestContextImpl) (interface{}, error) {
	b, err := ctx.ReadBodyBytes()
	if err != nil {
		log.Info("ERR recBlockVoteRound1Handler: ", err)
	}
	//log.Info("recBlockVoteRound1Handler: ", string(b))
	v := &Vote{}
	err = json.Unmarshal(b, v)
	if err != nil {
		log.Info("ERR recBlockVoteRound1Handler: ", err)
	}
	conn := redis.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("SADD", v.Hash+"round1", b)
	if err != nil {
		log.Info("recBlockVoteRound1Handler err SADD: ", err)
	}
	vc, err := redis.ToInt(conn.Do("SCARD", v.Hash+"round1")) // 获取投票个数
	if err != nil {
		log.Info("recBlockVoteRound1Handler err SCARD:", err)
	}
	log.Infof("收到第一轮投票 %d 张", vc)
	if vc == Nodes { // 这里必须收到所有人的票才能进入下一轮次，实际是收到2f+1张同意票即可
		//log.Info("voteForRoundTwo")
		go voteForRoundNew(v.Hash)
	}
	return nil, nil
}

// 查看同意票数是否达到2f+1
func voteForRoundNew(hash string) {
	//when receive whole nodes votes
	//then statistics
	conn := redis.Pool.Get()
	defer conn.Close()

	votes, err := rd.ByteSlices(conn.Do("SMEMBERS", hash+"round1"))
	if err != nil {
		log.Info("recBlockVoteRound1Handler err SADD: ", err)
	}
	var vs []Vote
	for _, data := range votes {
		t := Vote{}
		err := json.Unmarshal(data, &t)
		if err != nil {
			log.Info("recBlockVoteRound1Handler err Json:", err)
			continue
		}
		vs = append(vs, t)
	}
	votecount := 0
	for _, v := range vs {
		if v.Vote {
			votecount++
		}
	}
	v := false
	if float64(votecount) > float64(Nodes)*0.75 { // 实际无需 3/4 这么多
		log.Infof("第一轮投票同意票数 %d 张，达到2f+1，发起第二轮投票", votecount)
		v = true
	}
	rv := &ReVote{Sender: Sender, Vote: vs, Hash: hash, V: v}
	b, err := json.Marshal(rv)
	if err != nil {
		log.Info("voteForRoundTwo err json: ", err)
	}
	log.Info("第二轮投票, vote: ", v)
	Broadcast("recBlockVoteRound2", b)

}
