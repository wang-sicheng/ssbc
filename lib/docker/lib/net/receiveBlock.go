package net

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/cloudflare/cfssl/log"
	"github.com/ssbc/common"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
)

var (
	scStatus  map[string]string
	scVersion map[string]string
	scIP      map[string]string
	scHash    map[string]bool
)

type LcCommand struct {
	Command string
	SCNAME  string
	Version string
	Data    []byte
}

func receiveBlock(s *Server) *serverEndpoint {
	return &serverEndpoint{
		Methods: []string{"POST"},
		Handler: receiveBlockHandler,
		Server:  s,
	}
}

func receiveBlockHandler(ctx *serverRequestContextImpl) (interface{}, error) {
	b, err := ctx.ReadBodyBytes()
	if err != nil {
		log.Info("ERR receiveBlockHandler: ", err)
	}
	log.Info("receiveBlockHandler: ", string(b))
	newBlock := &common.Block{}
	err = json.Unmarshal(b, newBlock)
	if err != nil {
		log.Info("ERR receiveBlockHandler: ", err)
	}
	log.Info("receiveBlockHandler newBlock: ", newBlock)
	go smartContractExcu(newBlock)

	return nil, nil
}

func smartContractExcu(b *common.Block) {
	for _, tx := range b.TX {
		if IsSmartContract(&tx) {
			txToSmartContract(&tx)
		}
	}
}

func IsSmartContract(tx *common.Transaction) bool {
	if tx.To == "0" || tx.To == "lcssc" {
		return true
	}
	if _, ok := scStatus[tx.To]; ok {
		return true
	}
	return scHash[tx.To]
}

func txToSmartContract(tx *common.Transaction) {
	if tx.To == "lcssc" || tx.To == "0" {
		go lcssc(tx)
	} else {
		client := &Client{
			Url: scIP[tx.To],
		}
		client.httpClient = &http.Client{

			Transport: &http.Transport{

				MaxIdleConns:        10,
				MaxIdleConnsPerHost: 10,
				DisableKeepAlives:   false,
			},
		}
		endPoint := scIP[tx.To] + "/init"

		req, err := NewPost(endPoint, []byte(tx.Message))
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8;")

		err = client.SendReq(req, nil)
		if err != nil {
			return
		}
		newtx := common.Transaction{"lcssc", tx.From, time.Now().String(), "signature", "result"}
		sendTx(newtx)
	}
}

func lcssc(tx *common.Transaction) {
	lcc := &LcCommand{}
	err := json.Unmarshal([]byte(tx.Message), lcc)
	if err != nil {
		log.Info("lcc err json :", err)
		return
	}
	if tx.From == "lcssc" || tx.From == "0" {
		if lcc.Command == "install" {
			installSC(lcc)
		}
		if lcc.Command == "instantiate" {
			instantiateSC(lcc)
		}
	} else {
		switch lcc.Command {
		case "install":
			sci, _ := GenerateSCSpec(lcc, tx.From)
			scd, _ := Compile(sci)
			scHash[scd.Hash] = true
			b, _ := json.Marshal(scd)
			newlcc := &LcCommand{"install", scd.SCNAME, scd.Version, b}
			b, _ = json.Marshal(newlcc)
			newtx := common.Transaction{"lcssc", "lcssc", time.Now().String(), "signature", string(b)}
			sendTx(newtx)
		case "instantiate":
			if scStatus[lcc.SCNAME] != "install" {
				return
			}
			res, err := instantiate(lcc)
			if err != nil {
				return
			}
			scs := &SCInStatus{lcc.SCNAME, res}
			b, _ := json.Marshal(scs)
			newlcc := &LcCommand{"instantiate", lcc.SCNAME, lcc.Version, b}
			b, _ = json.Marshal(newlcc)
			newtx := common.Transaction{"lcssc", "lcssc", time.Now().String(), "signature", string(b)}
			sendTx(newtx)
		case "update":
			go update(lcc)
		}
	}
}

func installSC(lcc *LcCommand) {
	// docker install
	scStatus[lcc.SCNAME] = "install"
	scVersion[lcc.SCNAME] = lcc.Version

	//scIP[lcc.SCNAME] =
}

func instantiate(lcc *LcCommand) (string, error) {
	client := &Client{
		Url: scIP[lcc.SCNAME],
	}
	client.httpClient = &http.Client{

		Transport: &http.Transport{

			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			DisableKeepAlives:   false,
		},
	}
	endPoint := scIP[lcc.SCNAME] + "/init"

	req, err := NewPost(endPoint, lcc.Data)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8;")

	err = client.SendReq(req, nil)
	if err != nil {
		return "", err
	}
	return "result", nil
}

func instantiateSC(lcc *LcCommand) {
	scStatus[lcc.SCNAME] = "instantiate"
}

type SmartContractInit struct {
	SCNAME  string
	CREATER string
	Version string
	Data    []byte
}

type SmartContractDefinition struct {
	SCNAME  string
	Hash    string
	Version string
}

func GenerateSCSpec(lcc *LcCommand, creater string) (*SmartContractInit, error) {
	scinit := &SmartContractInit{}
	scinit.SCNAME = lcc.SCNAME
	scinit.Version = lcc.Version
	scinit.CREATER = creater
	scinit.Data = lcc.Data
	return scinit, nil
}

func Compile(sc *SmartContractInit) (*SmartContractDefinition, error) {
	if preCompole(sc) {
		hash, err := compile(sc)
		if err != nil {
			return nil, err
		}
		scd := &SmartContractDefinition{SCNAME: sc.SCNAME, Hash: hash, Version: sc.Version}
		return scd, nil
	}

	return nil, nil
}

func preCompole(sc *SmartContractInit) bool {
	// check if already exists
	//check version
	//check sig

	return true
}

func compile(sc *SmartContractInit) (string, error) {
	file, err := os.Create("D:/" + sc.SCNAME + ".go")
	if err != nil {
		log.Info("compile err file: ", err)
	}
	_, err = file.Write(sc.Data)
	if err != nil {
		log.Info("compile err write: ", err)
		return "", err
	}
	file.Close()
	cmd := exec.Command("go", "build", "-o", "D:/"+sc.SCNAME+".exe", "D:/"+sc.SCNAME+".go")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Error(err.Error(), stderr.String())
		return "", errors.New(stderr.String())
	}
	fp, err := ioutil.ReadFile("D:/" + sc.SCNAME + ".exe")
	if err != nil {
		log.Info("compile err openfile: ", err)
	}
	h := sha256.New()
	h.Write(fp)
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed), nil

}

func sendTx(tx common.Transaction) {
	b, _ := json.Marshal(tx)
	Broadcast("reciveTx", b)
}

type SCIns struct {
	SCNAME  string
	CREATER string
	Version string
	Data    []byte
}

type SCInStatus struct {
	SCNAME string
	Status string
}

type scupdate struct {
	SCName      string
	Version     string
	install     []byte
	instantiate []byte
}

func update(lcc *LcCommand) {
	//stop old version
	scu := &scupdate{}
	err := json.Unmarshal(lcc.Data, scu)
	if err != nil {
		return
	}
	installTx := common.Transaction{"", "lcssc", time.Now().String(), "signature", string(scu.install)}
	sendTx(installTx)
	for {
		if scStatus[lcc.SCNAME] == "install" {
			break
		}
		time.Sleep(time.Second)
	}
	instantiateTx := common.Transaction{"", "lcssc", time.Now().String(), "signature", string(scu.instantiate)}
	sendTx(instantiateTx)
}
