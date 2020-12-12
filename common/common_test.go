package common

import (
	"fmt"
	"strconv"
	"testing"
)

func Test_connect(t *testing.T) {
	f := MaxByzantiumNumber(10)
	q := QuorumNumber(10)
	fmt.Print("f: " + strconv.Itoa(f) + "\n")
	fmt.Print("q: " + strconv.Itoa(q) + "\n")
}
