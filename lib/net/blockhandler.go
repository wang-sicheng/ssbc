package net

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudflare/cfssl/log"
	rd "github.com/gomodule/redigo/redis"
	"github.com/ssbc/common"
	"github.com/ssbc/lib/redis"

	"bytes"
	"hash"

	"github.com/ssbc/crypto"
	"strconv"
)

type TestContent struct {
	x string
}

//CalculateHash hashes the values of a TestContent
func (t TestContent) CalculateHash() ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(t.x)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

//Equals tests for equality of two Contents
func (t TestContent) Equals(other Content) (bool, error) {
	return t.x == other.(TestContent).x, nil
}

type test struct {
	Name []string
}

func main() {
	//merkle tree tests
	//m:= []string{"e5","d4","c3","b2","a1"}
	//mm:=[][]byte{}
	//for _,t:=range m{
	//	//fmt.Println(t)
	//	mm = append(mm,[]byte(t))
	//}
	//mt := common.NewMerkleTree(mm)
	//fmt.Println(hex.EncodeToString(mt.RootNode.Data))
	//var list []Content
	//list = append(list, TestContent{x: "e5"})
	//list = append(list, TestContent{x: "d4"})
	//list = append(list, TestContent{x: "c3"})
	//list = append(list, TestContent{x: "b2"})
	//list = append(list, TestContent{x: "a1"})
	//a,_:= list[0].CalculateHash()
	//fmt.Println(a)
	//hash := sha256.Sum256([]byte("e5"))
	//fmt.Println(hash)
	////Create a new Merkle Tree from the list of Content
	//t, err := NewTree(list)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	////Get the Merkle Root of the tree
	//mr := t.MerkleRoot()
	//fmt.Println(hex.EncodeToString(mr))

	//
	//trans := make(chan []byte,100)
	//go marshalTrans(trans)
	//go transToRedis(trans)
	//i:=0
	//t1 := time.Now()
	//for{
	//
	//	i++
	//
	//	a:=pullTrans()
	//	t2 :=time.Now()
	//
	//	//close(trans)
	//	for _,t:=range a{
	//		tmp := common.Transaction{}
	//		err := json.Unmarshal(t,&tmp)
	//		if err!=nil{
	//			panic(err)
	//		}
	//		//log.Info("haha", tmp)
	//	}
	//	t3 :=time.Now()
	//
	//	log.Info(t2.Sub(t1))
	//	log.Info(t3.Sub(t2))
	//	if i>3{break}
	//	time.Sleep(time.Second*2)
	//}
	//conn := redis.Pool.Get()
	//	defer conn.Close()
	//	_,err := conn.Do("RPUSH", transjson...)
	//	if err != nil{
	//		log.Info("recTrans err: ", err)
	//	}

	//recTrans()
	//pullTrans()

}
func marshalTrans(trans chan []byte) {

	tx := generateTx()
	for _, data := range tx {
		if verifyTrans(data) {
			transbyte, jserr := json.Marshal(data)
			if jserr != nil {
				log.Info(jserr)
				return
			}

			trans <- transbyte

		}
	}

}

func recTrans() {
	//接受交易，验证 存入redis
	trans := generateTx()
	transjson := []interface{}{"transPool"}
	for _, data := range trans {
		if verifyTrans(data) {
			transbyte, jserr := json.Marshal(data)
			if jserr != nil {

				return
			}
			transjson = append(transjson, transbyte)

		}
	}
	conn := redis.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("RPUSH", transjson...)
	if err != nil {
		log.Info("recTrans err: ", err)
	}

}

func transToRedis(trans chan []byte) {
	conn := redis.Pool.Get()
	defer conn.Close()
	i := 1
	for t := range trans {
		i++
		_, err := conn.Do("RPUSH", "transPool", t)
		if err != nil {
			log.Info("recTrans err: ", err)
			return
		}
		log.Info("transToRedis", string(t))
		if i > 10000 {
			return
		}
	}
}

//TODO:交易验证
func verifyTrans(tran common.Transaction) bool {
	res := crypto.VerifySignECC([]byte(tran.Message), tran.Signature, tran.SenderPublicKey)
	return res
}

// 从redis中拉取交易
func pullTrans() [][]byte {
	conn := redis.Pool.Get()
	defer conn.Close()

	//tran,err := rd.ByteSlices(conn.Do("LRANGE","transPool",0,10))
	//if err != nil{
	//	log.Info("recTrans err lrange: ", err)
	//}
	//for _,tmp := range tran{
	//	t := common.Transaction{}
	//	err := json.Unmarshal(tmp, &t)
	//	if err!=nil{
	//		panic(err)
	//	}
	//	log.Info("tran : ", t)
	//}
	comTransSet := [][]byte{}
	//time.Sleep(time.Second*5)

	for {
		transs, err := rd.ByteSlices(conn.Do("BLPOP", "transPool", "1"))
		if err != nil {
			log.Info("BLPOP err: ", err)
		}
		if transs != nil {
			comTransSet = append(comTransSet, transs[1])
		}
		if transs == nil || len(comTransSet) >= transinblock {
			break
		}
	}
	return comTransSet

	//comTransSet := []interface{}{}
	//for{
	//	trans,err := conn.Do("BLPOP","transPool",1)
	//
	//	log.Info("err",err)
	//	log.Info("trans",trans)
	//	if trans != nil{
	//		comTransSet = append(comTransSet, trans)
	//	}
	//
	//	if trans == nil || len(comTransSet) > 100{
	//		break
	//	}
	//}
	//for i,tmp := range comTransSet{
	//	tran,err:= rd.ByteSlices(tmp,nil)
	//	if err != nil{
	//		log.Info(err)
	//	}
	//	t := common.Transaction{}
	//	err = json.Unmarshal(tran[0], &t)
	//	if err != nil{
	//		panic(err)
	//	}
	//	log.Info(i,t)
	//	log.Info(i,tran)
	//}

}

// 生成一系列交易
func generateTx() []common.Transaction {
	var res []common.Transaction
	var message string
	//message = "transaction message"
	//strSignature := crypto.SignECC([]byte(message), "eccprivate.pem")
	for i := 0; i <= 15000; i++ {
		message = "message"+strconv.Itoa(i)
		strSignature := crypto.SignECC([]byte(message), "eccprivate.pem")
		//cur := time.Now()
		tmp := common.Transaction{
			SenderAddress:   strconv.Itoa(i), //int(cur.Unix())+
			ReceiverAddress: "To",
			Timestamp:       strconv.Itoa(i), //cur.String(),
			Signature:       strSignature,
			Message:         message,
			SenderPublicKey: crypto.GetECCPublicKey("eccpublic.pem"),
		}
		res = append(res, tmp)
	}
	return res
}

type Content interface {
	CalculateHash() ([]byte, error)
	Equals(other Content) (bool, error)
}

//MerkleTree is the container for the tree. It holds a pointer to the root of the tree,
//a list of pointers to the leaf nodes, and the merkle root.
type MerkleTree struct {
	Root         *Node
	merkleRoot   []byte
	Leafs        []*Node
	hashStrategy func() hash.Hash
}

//Node represents a node, root, or leaf in the tree. It stores pointers to its immediate
//relationships, a hash, the content stored if it is a leaf, and other metadata.
type Node struct {
	Tree   *MerkleTree
	Parent *Node
	Left   *Node
	Right  *Node
	leaf   bool
	dup    bool
	Hash   []byte
	C      Content
}

//verifyNode walks down the tree until hitting a leaf, calculating the hash at each level
//and returning the resulting hash of Node n.
func (n *Node) verifyNode() ([]byte, error) {
	if n.leaf {
		return n.C.CalculateHash()
	}
	rightBytes, err := n.Right.verifyNode()
	if err != nil {
		return nil, err
	}

	leftBytes, err := n.Left.verifyNode()
	if err != nil {
		return nil, err
	}

	h := n.Tree.hashStrategy()
	if _, err := h.Write(append(leftBytes, rightBytes...)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

//calculateNodeHash is a helper function that calculates the hash of the node.
func (n *Node) calculateNodeHash() ([]byte, error) {
	if n.leaf {
		return n.C.CalculateHash()
	}

	h := n.Tree.hashStrategy()
	if _, err := h.Write(append(n.Left.Hash, n.Right.Hash...)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

//NewTree creates a new Merkle Tree using the content cs.
func NewTree(cs []Content) (*MerkleTree, error) {
	var defaultHashStrategy = sha256.New
	t := &MerkleTree{
		hashStrategy: defaultHashStrategy,
	}
	root, leafs, err := buildWithContent(cs, t)
	if err != nil {
		return nil, err
	}
	t.Root = root
	t.Leafs = leafs
	t.merkleRoot = root.Hash
	return t, nil
}

//NewTreeWithHashStrategy creates a new Merkle Tree using the content cs using the provided hash
//strategy. Note that the hash type used in the type that implements the Content interface must
//match the hash type profided to the tree.
func NewTreeWithHashStrategy(cs []Content, hashStrategy func() hash.Hash) (*MerkleTree, error) {
	t := &MerkleTree{
		hashStrategy: hashStrategy,
	}
	root, leafs, err := buildWithContent(cs, t)
	if err != nil {
		return nil, err
	}
	t.Root = root
	t.Leafs = leafs
	t.merkleRoot = root.Hash
	return t, nil
}

// GetMerklePath: Get Merkle path and indexes(left leaf or right leaf)
func (m *MerkleTree) GetMerklePath(content Content) ([][]byte, []int64, error) {
	for _, current := range m.Leafs {
		ok, err := current.C.Equals(content)
		if err != nil {
			return nil, nil, err
		}

		if ok {
			currentParent := current.Parent
			var merklePath [][]byte
			var index []int64
			for currentParent != nil {
				if bytes.Equal(currentParent.Left.Hash, current.Hash) {
					merklePath = append(merklePath, currentParent.Right.Hash)
					index = append(index, 1) // right leaf
				} else {
					merklePath = append(merklePath, currentParent.Left.Hash)
					index = append(index, 0) // left leaf
				}
				current = currentParent
				currentParent = currentParent.Parent
			}
			return merklePath, index, nil
		}
	}
	return nil, nil, nil
}

//buildWithContent is a helper function that for a given set of Contents, generates a
//corresponding tree and returns the root node, a list of leaf nodes, and a possible error.
//Returns an error if cs contains no Contents.
func buildWithContent(cs []Content, t *MerkleTree) (*Node, []*Node, error) {
	if len(cs) == 0 {
		return nil, nil, errors.New("error: cannot construct tree with no content")
	}
	var leafs []*Node
	for _, c := range cs {
		hash, err := c.CalculateHash()
		if err != nil {
			return nil, nil, err
		}

		leafs = append(leafs, &Node{
			Hash: hash,
			C:    c,
			leaf: true,
			Tree: t,
		})
	}
	if len(leafs)%2 == 1 {
		duplicate := &Node{
			Hash: leafs[len(leafs)-1].Hash,
			C:    leafs[len(leafs)-1].C,
			leaf: true,
			dup:  true,
			Tree: t,
		}
		leafs = append(leafs, duplicate)
	}
	root, err := buildIntermediate(leafs, t)
	if err != nil {
		return nil, nil, err
	}

	return root, leafs, nil
}

//buildIntermediate is a helper function that for a given list of leaf nodes, constructs
//the intermediate and root levels of the tree. Returns the resulting root node of the tree.
func buildIntermediate(nl []*Node, t *MerkleTree) (*Node, error) {
	var nodes []*Node
	for i := 0; i < len(nl); i += 2 {
		h := t.hashStrategy()
		var left, right int = i, i + 1
		if i+1 == len(nl) {
			right = i
		}
		chash := append(nl[left].Hash, nl[right].Hash...)
		if _, err := h.Write(chash); err != nil {
			return nil, err
		}
		n := &Node{
			Left:  nl[left],
			Right: nl[right],
			Hash:  h.Sum(nil),
			Tree:  t,
		}
		nodes = append(nodes, n)
		nl[left].Parent = n
		nl[right].Parent = n
		if len(nl) == 2 {
			return n, nil
		}
	}
	return buildIntermediate(nodes, t)
}

//MerkleRoot returns the unverified Merkle Root (hash of the root node) of the tree.
func (m *MerkleTree) MerkleRoot() []byte {
	return m.merkleRoot
}

//RebuildTree is a helper function that will rebuild the tree reusing only the content that
//it holds in the leaves.
func (m *MerkleTree) RebuildTree() error {
	var cs []Content
	for _, c := range m.Leafs {
		cs = append(cs, c.C)
	}
	root, leafs, err := buildWithContent(cs, m)
	if err != nil {
		return err
	}
	m.Root = root
	m.Leafs = leafs
	m.merkleRoot = root.Hash
	return nil
}

//RebuildTreeWith replaces the content of the tree and does a complete rebuild; while the root of
//the tree will be replaced the MerkleTree completely survives this operation. Returns an error if the
//list of content cs contains no entries.
func (m *MerkleTree) RebuildTreeWith(cs []Content) error {
	root, leafs, err := buildWithContent(cs, m)
	if err != nil {
		return err
	}
	m.Root = root
	m.Leafs = leafs
	m.merkleRoot = root.Hash
	return nil
}

//VerifyTree verify tree validates the hashes at each level of the tree and returns true if the
//resulting hash at the root of the tree matches the resulting root hash; returns false otherwise.
func (m *MerkleTree) VerifyTree() (bool, error) {
	calculatedMerkleRoot, err := m.Root.verifyNode()
	if err != nil {
		return false, err
	}

	if bytes.Compare(m.merkleRoot, calculatedMerkleRoot) == 0 {
		return true, nil
	}
	return false, nil
}

//VerifyContent indicates whether a given content is in the tree and the hashes are valid for that content.
//Returns true if the expected Merkle Root is equivalent to the Merkle root calculated on the critical path
//for a given content. Returns true if valid and false otherwise.
func (m *MerkleTree) VerifyContent(content Content) (bool, error) {
	for _, l := range m.Leafs {
		ok, err := l.C.Equals(content)
		if err != nil {
			return false, err
		}

		if ok {
			currentParent := l.Parent
			for currentParent != nil {
				h := m.hashStrategy()
				rightBytes, err := currentParent.Right.calculateNodeHash()
				if err != nil {
					return false, err
				}

				leftBytes, err := currentParent.Left.calculateNodeHash()
				if err != nil {
					return false, err
				}

				if _, err := h.Write(append(leftBytes, rightBytes...)); err != nil {
					return false, err
				}
				if bytes.Compare(h.Sum(nil), currentParent.Hash) != 0 {
					return false, nil
				}
				currentParent = currentParent.Parent
			}
			return true, nil
		}
	}
	return false, nil
}

//String returns a string representation of the node.
func (n *Node) String() string {
	return fmt.Sprintf("%t %t %v %s", n.leaf, n.dup, n.Hash, n.C)
}

//String returns a string representation of the tree. Only leaf nodes are included
//in the output.
func (m *MerkleTree) String() string {
	s := ""
	for _, l := range m.Leafs {
		s += fmt.Sprint(l)
		s += "\n"
	}
	return s
}
