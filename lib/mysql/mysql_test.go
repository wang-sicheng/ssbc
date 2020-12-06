package mysql

import (
	"testing"
)

func Test_connect(t *testing.T) {
	DB = InitDB()
	QueryAllBlocks(DB)
	CloseDB()
}
