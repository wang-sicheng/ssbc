package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"math/big"
	"os"
)

type ECDSASignature struct {
	R, S *big.Int
}
//生成ECC椭圆曲线密钥对，保存到文件
func GenerateECCKey() {
	//生成密钥对
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	//保存私钥
	//生成文件
	privatefile, err := os.Create("eccprivate.pem")
	if err != nil {
		panic(err)
	}
	//x509编码
	eccPrivateKey, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		panic(err)
	}
	//pem编码
	privateBlock := pem.Block{
		Type:  "ecc private key",
		Bytes: eccPrivateKey,
	}
	pem.Encode(privatefile, &privateBlock)
	//保存公钥
	publicKey := privateKey.PublicKey
	//创建文件
	publicfile, err := os.Create("eccpublic.pem")
	//x509编码
	eccPublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}
	//pem编码
	block := pem.Block{Type: "ecc public key", Bytes: eccPublicKey}
	pem.Encode(publicfile, &block)
}

//取得ECC私钥
func GetECCPrivateKey(path string) string {
	//读取私钥
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	return string(buf)
}

//取得ECC公钥
func GetECCPublicKey(path string) string {
	//读取公钥
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	publicKey := string(buf)
	return publicKey
}

//对消息生成数字签名
func SignECC(msg []byte, privateKeyStr string) string {
	//pem解码
	block, _ := pem.Decode([]byte(privateKeyStr))
	//x509解码
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	//计算哈希值
	hash := sha256.New()
	//填入数据
	hash.Write(msg)
	bytes := hash.Sum(nil)
	//对哈希值生成数字签名
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, bytes)
	if err != nil {
		panic(err)
	}
	//signatureAsn1, _ := asn1.Marshal(ECDSASignature{
	//	R: r,
	//	S: s,
	//})
	signatureJson, _ := json.Marshal(&ECDSASignature{
		R: r,
		S: s,
	})
	return string(signatureJson)
}

//验证数字签名
func VerifySignECC(msg []byte, signatureStr string, publicKeyStr string) bool {
	//pem解密
	block, _ := pem.Decode([]byte(publicKeyStr))
	//x509解密
	publicInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	publicKey := publicInterface.(*ecdsa.PublicKey)
	//计算哈希值
	hash := sha256.New()
	hash.Write(msg)
	bytes := hash.Sum(nil)
	////验证数字签名
	//signature := new(ECDSASignature)
	//_, err = asn1.Unmarshal([]byte(signatureAsn1), signature)
	//if err != nil {
	//	log.Error(signatureAsn1)
	//	log.Error(len(signatureAsn1))
	//	panic(err)
	//}
	signature := new(ECDSASignature)
	err = json.Unmarshal([]byte(signatureStr), signature)
	if err != nil {
		panic(err)
	}
	verify := ecdsa.Verify(publicKey, bytes, signature.R, signature.S)
	return verify
}

////测试
//func main() {
//	//生成ECC密钥对文件
//	//GenerateECCKey()
//	//for i := 0; i < 30000; i++ {
//	//	message := []byte(strconv.Itoa(i))
//	//	signature := SignECC(message, GetECCPrivateKey("eccprivate.pem"))
//	//	res := VerifySignECC(message, signature, GetECCPublicKey("eccpublic.pem"))
//	//	if res != true{
//	//		print(i)
//	//	}
//	//}
//	message := []byte("strconv.Itoa(i)")
//	signature := SignECC(message, GetECCPrivateKey("eccprivate.pem"))
//	res := VerifySignECC(message, signature, GetECCPublicKey("eccpublic.pem"))
//	fmt.Println(res)
//}
