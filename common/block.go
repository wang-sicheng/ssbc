package common

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type Block struct {
	Index     int `db:bIndex`
	Timestamp string `db:Timestamp`
	BPM       int `db:BPM`
	Hash      string `db:Hash`
	PrevHash  string `db:Prevhash`
	Merkle	  string
	TX 		  []Transaction
	Signature string
	// 每条交易要签名 是client的签名 tbd
	// merkel tree
	//  hash 块头hash
	//  块体 交易
}

type Transaction struct{
	From string
	To string
	Timestamp string
	Signature string
	Message string
}

var Blockchain []Block
var Blockchains = make(chan Block , 100000)

var Tx100 []Transaction

type Message struct {
	BPM int `BPM`
}

var mutex = &sync.Mutex{}

func run() error {
	mux := makeMuxRouter()
	httpPort := os.Getenv("PORT")
	log.Println("HTTP Server Listening on port :", httpPort)
	s := &http.Server{
		Addr:           ":" + httpPort,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// create handlers
func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	//muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

// write blockchain when we receive an http request
func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

// takes JSON payload as an input for heart rate (BPM)
//func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//	var msg Message
//
//	decoder := json.NewDecoder(r.Body)
//	if err := decoder.Decode(&msg); err != nil {
//		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
//		return
//	}
//	defer r.Body.Close()
//
//	mutex.Lock()
//	prevBlock := Blockchain[len(Blockchain)-1]
//	//newBlock := GenerateBlock(prevBlock, msg.BPM)
//
//	if isBlockValid(newBlock, prevBlock) {
//		Blockchain = append(Blockchain, newBlock)
//		spew.Dump(Blockchain)
//	}
//	mutex.Unlock()
//
//	respondWithJSON(w, r, http.StatusCreated, newBlock)
//
//}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

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
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.BPM) + block.PrevHash
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
	newBlock.Hash = calculateHash(newBlock)
	newBlock.Merkle = "Merkle"
	newBlock.Signature = "Signature"
	newBlock.TX = Tx100
	return newBlock
}

func Init(){
	Tx100 = generateTx()
	b := Block{0,"0",0,"0","0","0",nil,"0"}
	Blockchain = append(Blockchain, b)
	Blockchains <- b

}

func generateTx()[]Transaction{
	res := []Transaction{}
	for i := 0; i < 10; i++{
		tmp := Transaction{
			From:"From",
			To:"To",
			Timestamp:"Timestamp",
			Signature:"Signature",
			Message:"Message",
		}
		res = append(res, tmp)
	}
	return res
}
