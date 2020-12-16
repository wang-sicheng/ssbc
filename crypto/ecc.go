package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
)

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
func GetECCPrivateKey(path string) *ecdsa.PrivateKey {
	//读取私钥
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	//pem解码
	block, _ := pem.Decode(buf)
	//x509解码
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	return privateKey
}

////取得ECC公钥
//func GetECCPublicKey(path string) *ecdsa.PublicKey {
//	//读取公钥
//	file, err := os.Open(path)
//	if err != nil {
//		panic(err)
//	}
//	defer file.Close()
//	info, _ := file.Stat()
//	buf := make([]byte, info.Size())
//	file.Read(buf)
//	//pem解密
//	block, _ := pem.Decode(buf)
//	//x509解密
//	publicInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
//	if err != nil {
//		panic(err)
//	}
//	publicKey := publicInterface.(*ecdsa.PublicKey)
//	return publicKey
//}

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

//对消息的散列值生成数字签名
func SignECC(msg []byte, path string) string {
	//取得私钥
	privateKey := GetECCPrivateKey(path)
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
	//rtext, _ := r.MarshalText()
	//stext, _ := s.MarshalText()
	//- 将r、s转成r、s字符串
	strSigR := fmt.Sprintf("%x", r)
	strSigS := fmt.Sprintf("%x", s)
	//fmt.Printf("r的16进制为:%s,长度为%d\n", strSigR, len(strSigR))
	//fmt.Printf("s的16进制为:%s,长度为%d\n", strSigS, len(strSigS))
	if len(strSigR) == 63 {
		strSigR = "0" + strSigR
	}
	if len(strSigS) == 63 {
		strSigR = "0" + strSigS
	}
	//形成数字签名的der格式
	derString := MakeDerSign(strSigR, strSigS)
	return derString
}

////验证数字签名
//func VerifySignECC(msg []byte, derSignString string, path string) bool {
//	//读取公钥
//	publicKey := GetECCPublicKey(path)
//	//计算哈希值
//	hash := sha256.New()
//	hash.Write(msg)
//	bytes := hash.Sum(nil)
//	//验证数字签名
//	rBytes, sBytes := ParseDERSignString(derSignString)
//	r := new(big.Int).SetBytes(rBytes)
//	s := new(big.Int).SetBytes(sBytes)
//	verify := ecdsa.Verify(publicKey, bytes, r, s)
//	return verify
//}

//验证数字签名
func VerifySignECC(msg []byte, derSignString string, publicKeyStr string) bool {
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
	//验证数字签名
	rBytes, sBytes := ParseDERSignString(derSignString)
	r := new(big.Int).SetBytes(rBytes)
	s := new(big.Int).SetBytes(sBytes)
	verify := ecdsa.Verify(publicKey, bytes, r, s)
	return verify
}

//生成数字签名的DER编码格式
func MakeDerSign(strR, strS string) string {
	//获取R和S的长度
	lenSigR := len(strR) / 2 //16进制每两位1字节
	lenSigS := len(strS) / 2
	//- 计算DER序列的总长度
	l := lenSigR + lenSigS + 4
	//- 将10进制长度转16进制字符串
	strLenSigR := fmt.Sprintf("%x", int64(lenSigR))
	strLenSigS := fmt.Sprintf("%x", int64(lenSigS))
	strLen := fmt.Sprintf("%x", int64(l))
	//- 拼凑DER编码格式
	derString := "30" + strLen
	derString += "02" + strLenSigR + strR
	derString += "02" + strLenSigS + strS
	derString += "01"
	return derString
}

// 解析DER编码格式
func ParseDERSignString(derString string) (rBytes, sBytes []byte) {
	derBytes, _ := hex.DecodeString(derString)
	rBytes = derBytes[4:36]
	sBytes = derBytes[len(derBytes)-33 : len(derBytes)-1]
	//strSigR := fmt.Sprintf("%x", rBytes)
	//strSigS := fmt.Sprintf("%x", sBytes)
	//fmt.Printf("rBytes的16进制为:%s,长度为%d\n", strSigR, len(strSigR))
	//fmt.Printf("sBytes的16进制为:%s,长度为%d\n", strSigS, len(strSigS))
	return rBytes, sBytes
}

////测试
//func main() {
//	//生成ECC密钥对文件
//	GenerateECCKey()
//
//	//模拟发送者
//	//要发送的消息
//	msg:=[]byte("hello world")
//	//生成数字签名
//	signature:=SignECC(msg,"eccprivate.pem")
//
//	//模拟接受者
//	//接受到的消息
//	acceptmsg:=[]byte("hello world")
//	//接收到的签名
//	acceptSignature :=signature
//	//验证签名
//	verifySignECC := VerifySignECC(acceptmsg, acceptSignature, GetECCPublicKey("eccpublic.pem"))
//	fmt.Println("验证结果：",verifySignECC)
//}
