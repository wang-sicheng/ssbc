package net

import (
	"encoding/json"
	"github.com/cloudflare/cfssl/log"
	"time"
)

type Sig struct{
	Keys []byte
	Sender string
}

func syncPublickey(s *Server)*serverEndpoint{
	return &serverEndpoint{
		Methods: []string{ "POST"},
		Handler: syncPublickeyHandler,
		Server:  s,
	}
}

func syncPublickeyHandler(ctx *serverRequestContextImpl) (interface{}, error) {
	b,err := ctx.ReadBodyBytes()
	if err != nil{
		log.Info("syncPlulickey Readbody err: ", err)
		return nil, err
	}
	sig := Sig{}
	err = json.Unmarshal(b, &sig)
	if err != nil{
		log.Info("syncPlulickey json err: ", err)
		return nil, err
	}
	log.Info("syncPlulickey receiving: ", string(b))
	signatures[sig.Sender] = sig.Keys
	senders[ctx.req.RemoteAddr] = sig.Sender
	return nil, nil

}

func SendPublicKey(){


	for len(signatures) < Nodes{
		log.Info("Waiting for other nodes signature")
		time.Sleep(time.Second)
	}
	log.Info("Signature synchronize successfully")
	//do next
}