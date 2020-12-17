package net

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/ssbc/crypto"
	"github.com/ssbc/common"
	_ "github.com/ssbc/lib/mysql"
	"golang.org/x/crypto/ripemd160"
	"github.com/cloudflare/cfssl/log"
)

const version = byte(0x00)
const addressChecksumLen = 4 //对校验位一般取4位

func newAccount(s *Server) *serverEndpoint {
	return &serverEndpoint{
		Methods:   []string{"GET", "POST", "HEAD"},
		Handler:   newAccountHandler,
		Server:    s,
		successRC: 200,
	}
}

func newAccountHandler(ctx *serverRequestContextImpl) (interface{}, error) {
	b, err := ctx.ReadBodyBytes()
	if err != nil {
		log.Info("ERR newAccountHandler: ", err)
	}
	log.Info("newAccountHandler rec: ", string(b))
	ac := common.Account{}
	privateKey, publicKey := GenerateECCKey()
	xPrivateKey, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	//pem编码
	privateBlock := pem.Block{
		Bytes: xPrivateKey,
	}
	xPublicKey, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}
	publicBlock := pem.Block{
		Bytes: xPublicKey,
	}
	// 生成账户地址
	addr := GetAddress(publicKey)
	ac.Address = fmt.Sprintf("%s", addr)
	ac.PublicKey = fmt.Sprintf("%s", pem.EncodeToMemory(&publicBlock))
	ac.PrivateKey = fmt.Sprintf("%s", pem.EncodeToMemory(&privateBlock))
	//mysql.InsertAccount(ac)  // 存入数据库
	return ac.Address, nil
}

//生成ECC椭圆曲线密钥对
func GenerateECCKey() (*ecdsa.PrivateKey, []byte){
	//生成密钥对
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	//生成公钥
	publicKey := append(privateKey.PublicKey.X.Bytes(),privateKey.PublicKey.Y.Bytes()...)
	return privateKey, publicKey
}


//生成一个账户地址
func GetAddress(publicKey []byte) []byte {
	//调用公钥哈希函数，实现RIPEMD160(SHA256(Public Key))
	pubKeyHash := HashPubKey(publicKey)
	//存储version和公钥哈希的切片
	versionedPayload := append([]byte{version},pubKeyHash...)
	//调用checksum函数，对上面的切片进行双重哈希后，取出哈希后的切片的前面部分作为检验位的值
	checksum := checksum(versionedPayload)
	//把校验位加到上面切片后面
	fullPayload := append(versionedPayload,checksum...)
	//通过base58编码上述切片得到地址
	address := crypto.Base58Encode(fullPayload)
	return address
}


//公钥哈希函数，实现RIPEMD160(SHA256(Public Key))
func HashPubKey(pubKey []byte) []byte {
	//先hash公钥
	publicSHA256 := sha256.Sum256(pubKey)
	//对公钥哈希值做 ripemd160运算
	RIPEMD160Hasher := ripemd160.New()
	_,err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Info("ripemd160 error: ", err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}


//校验位checksum,双重哈希运算
func checksum(payload []byte) []byte {
	//下面双重哈希payload，在调用中，所引用的payload为（version + Pub Key Hash）
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	//addressChecksumLen代表保留校验位长度
	return secondSHA[:addressChecksumLen]
}

//func main() {
//	res, err := newAccountHandler()
//	fmt.Println(res, err)
//}

