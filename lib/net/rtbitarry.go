package net

import (
	"encoding/json"
	"github.com/cloudflare/cfssl/log"
	"github.com/ssbc/common"
	"github.com/ssbc/lib/redis"
	rd "github.com/gomodule/redigo/redis"
)




type TransHash struct {
	BlockHash string
	TransHashs []string
}

func receive_trans_bitarry(s *Server)*serverEndpoint{
	return &serverEndpoint{
		Methods: []string{ "POST"},
		Handler: rtbitarryHandler,
		Server:  s,
	}
}

//接收交易hash
func rtbitarryHandler(ctx *serverRequestContextImpl) (interface{}, error) {
	b,err := ctx.ReadBodyBytes()
	if err != nil{
		log.Info("rtbitarry Readbody err: ", err)
		return nil, err
	}
	transHash := TransHash{}
	err = json.Unmarshal(b, &transHash)
	if err != nil{
		log.Info("rtbitarry json err: ", err)
		return nil, err
	}
	log.Info("rtbitarryHandler receiving: ", string(b))
	go findCommonTrans(transHash, ctx.req.RemoteAddr)
	return nil, nil
}

func findCommonTrans(trans TransHash, sender string){
	
	currentBlock := blockState.GetCurrB()
	if trans.BlockHash != currentBlock.Hash{
		log.Info("findCommonTrans: BlockHash mismatch. This round may finish.")
		return
	}
	conn := redis.Pool.Get()
	defer conn.Close()
	data := []interface{}{"CommonTx"+ currentBlock.Hash + sender}
	for _,d := range trans.TransHashs{
		data = append(data, d)
	}
	//redis the trans
	_,err := conn.Do("SADD", data...)
	if err != nil{
		log.Info("findCommonTrans err SADD trans: ", err)
	}
	//redis the senders
	_,err = conn.Do("SADD", "CommonTx"+ currentBlock.Hash, "CommonTx"+ currentBlock.Hash+ sender)
	if err != nil{
		log.Info("findCommonTrans err SADD senders: ", err)
	}
	//check the len of the senders
	l,err := rd.Int(conn.Do("SCARD", "CommonTx"+ currentBlock.Hash))
	if err !=nil{
		log.Info("findCommonTrans err SCARD: ", err)
	}
	//Leader mode and check if got enough nodes tranx
	if !isSelfLeader {
		log.Info("findCommonTrans: Not Leader", isSelfLeader)
		return
	}
	log.Info("Leader Mode")
	if l != Nodes{
		log.Info("findCommonTrans: Do not get enough nodes ", l)
		return
	}

	//find the common trans
	senders,err := rd.Strings(conn.Do("SMEMBERS", "CommonTx"+ currentBlock.Hash))
	if err !=nil{
		log.Info("findCommonTrans err SMEMBERS: ", err)
	}
	senderInterface := []interface{}{}
	for _,s := range senders{
		senderInterface = append(senderInterface, s)
	}
	log.Info("senderInterface: ", senderInterface)
	commonTrans,err := rd.Strings(conn.Do("SINTER", senderInterface...))
	if err !=nil{
		log.Info("findCommonTrans err SINTER: ", err)
	}
	log.Info("findCommonTrans commonTrans: ", commonTrans)
	generateFromCommonTx(commonTrans, currentBlock)
}

func generateFromCommonTx(commonTrans []string, currentBlock common.Block){
	conn := redis.Pool.Get()
	defer conn.Close()
	mb,err := rd.Bytes(conn.Do("GET", "CommonTxCache"+ currentBlock.Hash))
	if err !=nil{
		log.Info("generateFromCommonTx err GET: ", err)
	}
	m := make(map[string][]byte)
	err = json.Unmarshal(mb, &m)
	if err !=nil{
		log.Info("generateFromCommonTx err GET: ", err)
	}
	trans := []common.Transaction{}
	for _,s := range commonTrans{
		if v,ok := m[s]; ok{
			tx := common.Transaction{}
			err := json.Unmarshal(v, &tx)
			if err != nil{
				log.Info("generateFromCommonTx err json trans: ", err)
				continue
			}
			trans = append(trans, tx)
		}
	}
	b := common.Block{TX:trans}
	b = common.GenerateBlock(currentBlock, b)
	bb, err := json.Marshal(b)
	if err != nil{
		log.Info("generateFromCommonTx err json block: ", err)
		return
	}
	Broadcast("recBlock",bb)
}