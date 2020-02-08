package common

import (
	"crypto/sha256"
)

// MerkleTree represent a Merkle tree
type MerkleTree struct {
	RootNode *MerkleNode
}

// MerkleNode represent a Merkle tree node
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

// NewMerkleTree creates a new Merkle tree from a sequence of data
func NewMerkleTree(data [][]byte) *MerkleTree {
	if len(data) == 0{
		return &MerkleTree{NewMerkleNode(nil,nil,nil)}
	}

	var nodes []MerkleNode



	for _, datum := range data {
		node := NewMerkleNode(nil, nil, datum)
		nodes = append(nodes, *node)
	}

	for len(nodes) != 1 {
		var newLevel []MerkleNode
		var node MerkleNode
		for j := 0; j < len(nodes); j += 2 {

			if j+1 == len(nodes){
				node = *NewMerkleNode(&nodes[j], &nodes[j], nil)
			}else{
				node = *NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			}
			//node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			newLevel = append(newLevel, node)
		}

		nodes = newLevel
	}

	mTree := MerkleTree{&nodes[0]}

	return &mTree
}

// NewMerkleNode creates a new Merkle tree node
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	mNode := MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		mNode.Data = hash[:]
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		mNode.Data = hash[:]
	}

	mNode.Left = left
	mNode.Right = right

	return &mNode
}