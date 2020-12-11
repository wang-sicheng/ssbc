package mysql

import (
	"github.com/ssbc/common"
	"testing"
)

func Test_connect(t *testing.T) {

	InitDB()
	var block common.Block
	block.Signature = "1"
	block.MerkleRoot = "fff"
	block.Hash = "hash"
	block.PrevHash = "pre"
	InsertBlock(block)
}
