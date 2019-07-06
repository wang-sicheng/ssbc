package net

import (
	"github.com/ssbc/common"
	"github.com/ssbc/lib/mysql"
)

type TestInfo struct {

	TName string

	Version string
}

func newTestInfoEndpoint(s *Server) *serverEndpoint {
	return &serverEndpoint{
		Methods: []string{"GET", "POST", "HEAD"},
		Handler: testinfoHandler,
		Server:  s,
	}
}

type TestInfoResponseNet struct {
	TName string

	Version string
}
func testinfoHandler(ctx *serverRequestContextImpl) (interface{}, error) {
	b := common.Block{5,"555",5,"h","pre"}
	mysql.InsertBlock(b)
	resp := &TestInfoResponseNet{}
	resp.TName = "hello world"
	resp.Version = "SSBC v0.1"
	return resp, nil
}