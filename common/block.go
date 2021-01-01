package common

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/cloudflare/cfssl/log"
	"github.com/ssbc/crypto"
	"strconv"
	"sync"
	"time"
)

type Block struct {
	Id         int    `json:"id"`
	Pid        int    `json:"pid"`
	PrevHash   string `json:"prev_hash"`
	Hash       string `json:"hash"`
	MerkleRoot string `json:"merkle_root"`
	TxCount    int    `json:"tx_count"`
	Signature  string `json:"signature"`
	Timestamp  string `json:"timestamp"`
	TX         []Transaction
}

type Transaction struct {
	Id              int    `json:"id"`
	BlockId         int    `json:"block_id"`
	SenderAddress   string `json:"sender_address"`
	ReceiverAddress string `json:"receiver_address"`
	Timestamp       string `json:"timestamp"`
	Signature       string `json:"signature"`
	Message         string `json:"message"`
	SenderPublicKey string `json:"sender_public_key"`
	TransferAmount  int    `json:"transfer_amount"`
}

type Account struct {
	Address string `json:"address"`
	PrivateKey string `json:"private_key"`
	PublicKey string `json:"public_key"`
}

//var Blockchains = make(chan Block, 100000)

var B Block

var Tx100 []Transaction

var mutex = &sync.Mutex{}

// make sure block is valid by checking index, and comparing the hash of the previous block
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Id+1 != newBlock.Id {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

//定义hash算法 SHA256 hasing
func calculateHash(block Block) string {
	//strconv.Itoa 将整数转换为十进制字符串形式（即：FormatInt(i, 10) 的简写）
	record := strconv.Itoa(block.Id) + block.Timestamp + block.PrevHash + block.MerkleRoot
	h := sha256.New()       //创建一个Hash对象
	h.Write([]byte(record)) //h.Write写入需要哈希的内容
	hashed := h.Sum(nil)    //h.Sum添加额外的[]byte到当前的哈希中，一般不是经常需要这个操作
	return hex.EncodeToString(hashed)
}

// create a new block using previous block's hash
func GenerateBlock(oldBlock Block, newBlock Block) Block {

	t := time.Now()
	newBlock.Id = oldBlock.Id + 1
	newBlock.Timestamp = t.String()
	newBlock.PrevHash = oldBlock.Hash
	newBlock.MerkleRoot = newBlock.GenerateMerkelRoot()
	newBlock.Hash = calculateHash(newBlock)
	strSignature := crypto.SignECC([]byte(newBlock.Hash), crypto.GetECCPrivateKey("eccprivate.pem"))
	newBlock.Signature = strSignature
	return newBlock
}

func Init() {
	Tx100 = generateTx()
	genesisBlock := Block{1, 0, "", "", "", 44, "", "", nil}
	genesisBlock.Hash = calculateHash(genesisBlock)
	genesisBlock.MerkleRoot = genesisBlock.GenerateMerkelRoot()
	log.Info("创世区块 id: ", genesisBlock.Id)
	//Blockchains <- genesisBlock
	B = genesisBlock

}

func generateTx() []Transaction {
	res := []Transaction{}
	for i := 0; i < 100; i++ {
		curTime := time.Now()
		tmp := Transaction{
			SenderAddress:   strconv.Itoa(curTime.Second()),
			ReceiverAddress: "To",
			Timestamp:       curTime.String(),
			Signature:       "Signature",
			Message:         "Message",
		}
		res = append(res, tmp)
	}
	return res
}

func (b *Block) GenerateMerkelRoot() string {
	mt := NewMerkleTree(transToByte(b.TX))
	return hex.EncodeToString(mt.RootNode.Data)
}

func transToByte(trans []Transaction) [][]byte {
	res := [][]byte{}
	for _, data := range trans {
		res = append(res, transTobyte(data))
	}
	return res
}
func transTobyte(tran Transaction) []byte {
	tranString := tran.SenderAddress + tran.ReceiverAddress + tran.Timestamp + tran.Signature + tran.Message
	return []byte(tranString)
}

func TransToByte(trans []Transaction) [][]byte {
	res := [][]byte{}
	for _, data := range trans {
		res = append(res, transTobyte(data))
	}
	return res
}
