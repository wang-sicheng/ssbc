package net

import (
	"encoding/json"
	"github.com/cloudflare/cfssl/log"
	rd "github.com/gomodule/redigo/redis"
	"github.com/ssbc/common"
	"github.com/ssbc/lib/redis"
)

type TransHash struct {
	BlockHash  string
	TransHashs []string
}

func receive_trans_bitarry(s *Server) *serverEndpoint {
	return &serverEndpoint{
		Methods: []string{"POST"},
		Handler: rtbitarryHandler,
		Server:  s,
	}
}

//接收交易hash
func rtbitarryHandler(ctx *serverRequestContextImpl) (interface{}, error) {
	b, err := ctx.ReadBodyBytes()
	if err != nil {
		log.Info("rtbitarry Readbody err: ", err)
		return nil, err
	}
	transHash := TransHash{}
	err = json.Unmarshal(b, &transHash)
	if err != nil {
		log.Info("rtbitarry json err: ", err)
		return nil, err
	}
	log.Infof("接收到 %d 条交易的Hash", len(transHash.TransHashs))
	//log.Info("rtbitarryHandler receiving: ", string(b))
	if IsSelfLeader {
		//go findCommonTrans(transHash, ctx.req.RemoteAddr)
		findCommonTrans(transHash, ctx.req.RemoteAddr)
	}
	return nil, nil
}

func findCommonTrans(trans TransHash, sender string) {

	currentBlock := blockState.GetCurrB()
	if trans.BlockHash != currentBlock.Hash { // recTransHash接收的body内容：1.上一个Block的Hash；2.收到的交易的Hash
		log.Info("findCommonTrans: BlockHash mismatch. This round may finish.")
		return
	}
	conn := redis.Pool.Get()
	defer conn.Close()
	data := []interface{}{"CommonTx" + currentBlock.Hash + sender}
	for _, d := range trans.TransHashs {
		data = append(data, d)
	}
	//redis the trans
	_, err := conn.Do("SADD", data...) // 将所有的交易Hash存入集合中
	if err != nil {
		log.Info("findCommonTrans err SADD trans: ", err)
	}
	//redis the senders									// 将sender存入集合中
	_, err = conn.Do("SADD", "CommonTx"+currentBlock.Hash, "CommonTx"+currentBlock.Hash+sender)
	if err != nil {
		log.Info("findCommonTrans err SADD senders: ", err)
	}
	//check the len of the senders
	l, err := rd.Int(conn.Do("SCARD", "CommonTx"+currentBlock.Hash))
	if err != nil {
		log.Info("findCommonTrans err SCARD: ", err)
	}
	//Leader mode and check if got enough nodes tranx
	if !IsSelfLeader {
		log.Info("findCommonTrans: Not Leader", IsSelfLeader)
		return
	}
	if l != Nodes {
		log.Info("收到交易集: ", l)
		return
	}

	//find the common trans
	senders, err := rd.Strings(conn.Do("SMEMBERS", "CommonTx"+currentBlock.Hash))
	if err != nil {
		log.Info("findCommonTrans err SMEMBERS: ", err)
	}
	senderInterface := []interface{}{}
	for _, s := range senders {
		senderInterface = append(senderInterface, s)
	}
	//log.Info("senderInterface: ", senderInterface)
	commonTrans, err := rd.Strings(conn.Do("SINTER", senderInterface...)) // SINTER返回指定集合的交集
	if err != nil {
		log.Info("findCommonTrans err SINTER: ", err)
	}
	log.Info("redis交易公共集长度: ", len(commonTrans))
	generateFromCommonTx(commonTrans, currentBlock)
}

func generateFromCommonTx(commonTrans []string, currentBlock common.Block) {
	conn := redis.Pool.Get()
	defer conn.Close()
	mb, err := rd.Bytes(conn.Do("GET", "CommonTxCache"+currentBlock.Hash)) // 每个节点在广播交易Hash时，会把交易缓存起来
	if err != nil {
		log.Info("generateFromCommonTx err GET: ", err)
	}
	m := make(map[string][]byte)
	err = json.Unmarshal(mb, &m)
	if err != nil {
		log.Info("generateFromCommonTx err GET: ", err)
	}
	trans := []common.Transaction{}
	for _, s := range commonTrans {
		if v, ok := m[s]; ok {
			tx := common.Transaction{}
			err := json.Unmarshal(v, &tx)
			if err != nil {
				log.Info("generateFromCommonTx err json trans: ", err)
				continue
			}
			trans = append(trans, tx)
		}
	}
	b := common.Block{TX: trans}
	b = common.GenerateBlock(currentBlock, b) // 建块
	bb, err := json.Marshal(b)
	if err != nil {
		log.Info("generateFromCommonTx err json block: ", err)
		return
	}
	log.Infof("广播id: %d 的区块", b.Id)
	Broadcast("recBlock", bb)
}
