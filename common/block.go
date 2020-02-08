package common

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"sync"
	"time"
	"github.com/cloudflare/cfssl/log"
)

type Block struct {
	Index     int `db:bIndex`
	Timestamp string `db:Timestamp`
	PrevHash  string `db:Prevhash`
	Merkle	  string `db:Merkle`
	Signature string `db:Signature`
	Hash      string `db:Hash`
	TX 		  []Transaction
}

type BlockHeader struct{
	Index     int `db:bIndex`
	Timestamp string `db:Timestamp`
	BPM       int `db:BPM`
	Hash      string `db:Hash`
	PrevHash  string `db:Prevhash`
	Merkle	  string
}

type Transaction struct{
	From string
	To string
	Timestamp string
	Signature string
	Message string
}

var Blockchains = make(chan Block , 100000)

var Tx100 []Transaction

var mutex = &sync.Mutex{}



// make sure block is valid by checking index, and comparing the hash of the previous block
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
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

// SHA256 hasing
func calculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + block.PrevHash + block.Merkle + block.Signature
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// create a new block using previous block's hash
func GenerateBlock(oldBlock Block, newBlock Block) Block {

	t := time.Now()
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Merkle = newBlock.GenerateMerkelRoot()
	newBlock.Signature = "Signature"
	newBlock.Hash = calculateHash(newBlock)
	return newBlock
}

func Init(){
	Tx100 = generateTx()
	genesisBlock  := Block{0,"","","","","",nil}
	genesisBlock .Hash = calculateHash(genesisBlock)
	genesisBlock.Merkle = genesisBlock.GenerateMerkelRoot()
	log.Info("GenesisBlock: ",genesisBlock)
	Blockchains <- genesisBlock

}

func generateTx()[]Transaction{
	res := []Transaction{}
	for i := 0; i < 100; i++{
		curTime := time.Now()
		tmp := Transaction{
			From:strconv.Itoa(curTime.Second()),
			To:"To",
			Timestamp:curTime.String(),
			Signature:"Signature",
			Message:"Message",
		}
		res = append(res, tmp)
	}
	return res
}

func (b *Block) GenerateMerkelRoot() string{
	mt := NewMerkleTree(transToByte(b.TX))
	return hex.EncodeToString(mt.RootNode.Data)
}

func transToByte(trans []Transaction)[][]byte{
	res := [][]byte{}
	for _,data := range trans{
		res = append(res, transTobyte(data))
	}
	return res
}
func transTobyte(tran Transaction)[]byte{
	tranString := tran.From + tran.To + tran.Timestamp + tran.Signature + tran.Message
	return []byte(tranString)
}

func TransToByte(trans []Transaction)[][]byte{
	res := [][]byte{}
	for _,data := range trans{
		res = append(res, transTobyte(data))
	}
	return res
}