package docker

import (
	"fmt"
	"os/exec"
	"github.com/ssbc/common"
	"encoding/json"
	"github.com/cloudflare/cfssl/log"
	"bytes"
	"os"
	"errors"
	"io/ioutil"
	"encoding/hex"
	"crypto/sha256"
	"github.com/spf13/viper"
	docker "github.com/fsouza/go-dockerclient"
)
func main() {
	file,err := os.Create("D:/hello.go")
	if err != nil{
		fmt.Println(err)
	}

	code := []byte(
`package main

import (
	"fmt"
	
)

func main() {
	fmt.Println("Hello World")
}
	`)
	_,err = file.Write(code)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("write succeff")
	file.Close()
	cmd := exec.Command("go", "build", "D:/hello.go")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(stderr.String())
		log.Info("here1")
		log.Error(err.Error(), stderr.String())
	} else {
		log.Info("here2")
		log.Info(out.String())
	}

}

type SmartContractInit struct{
	SCNAME string
	CREATER string
	Version string
	Data []byte
}

type SmartContractDefinition struct{
	SCNAME string
	Hash string
	Version string

}

func IsSmartContract(tx *common.Transaction)bool{
	if  tx.To == "0"{
		return true
	}
	return false
}

func GenerateSCSpec(tx *common.Transaction)(*SmartContractInit, error){
	scinit := &SmartContractInit{}
	err := json.Unmarshal([]byte(tx.Message), scinit)
	if err != nil{
		log.Info("generateSCSpec err json :", err)
		return nil, err
	}
	return scinit, nil
}

func transToDocker(){

}

func Compile(sc *SmartContractInit)(*SmartContractDefinition,error){
	if preCompole(sc){
		hash,err := compile(sc)
		if err != nil{
			return nil, err
		}
		scd := & SmartContractDefinition{SCNAME: sc.SCNAME, Hash:hash, Version: sc.Version}
		return scd, nil
	}



	return nil, nil
}

func preCompole(sc *SmartContractInit)bool{
	// check if already exists
	//check version
	//check sig

	return true
}

func compile(sc *SmartContractInit)(string, error){
	file,err := os.Create("D:/"+ sc.SCNAME+ ".go")
	if err != nil{
		log.Info("compile err file: ", err)
	}
	_,err = file.Write(sc.Data)
	if err != nil {
		log.Info("compile err write: ", err)
		return "", err
	}
	file.Close()
	cmd := exec.Command("go", "build", "-o", "D:/"+sc.SCNAME+ ".exe", "D:/"+ sc.SCNAME+ ".go")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Error(err.Error(), stderr.String())
		return "", errors.New(stderr.String())
	}
	fp,err := ioutil.ReadFile("D:/"+ sc.SCNAME+ ".exe")
	if err != nil{
		log.Info("compile err openfile: ", err)
	}
	h := sha256.New()
	h.Write(fp)
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed), nil

}

func execute(){

}

func NewDockerClient() (client *docker.Client, err error) {
	endpoint := viper.GetString("vm.endpoint")
	client, err = docker.NewClient(endpoint)

	return
}