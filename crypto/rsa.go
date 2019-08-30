package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/pem"
	"crypto/x509"
	"fmt"
	"os"
)

func GenRsaKey(bits int) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	privatefile, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	defer  privatefile.Close()
	err = pem.Encode(privatefile, block)
	if err != nil {
		return err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	publicfile, err := os.Create("public.pem")
	if err != nil {
		return err
	}
	defer publicfile.Close()
	err = pem.Encode(publicfile, block)
	if err != nil {
		return err
	}
	return nil
}

//RSA加密
func RSA_Encrypt(plainText []byte,path string)[]byte{
	//打开文件
	file,err:=os.Open(path)
	if err!=nil{
		panic(err)
	}
	defer file.Close()
	//读取文件的内容
	info, _ := file.Stat()
	buf:=make([]byte,info.Size())
	file.Read(buf)
	//pem解码
	block, _ := pem.Decode(buf)
	//x509解码

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err!=nil{
		panic(err)
	}

	//类型断言
	publicKey:=publicKeyInterface.(*rsa.PublicKey)
	//对明文进行加密
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err!=nil{
		panic(err)
	}
	//返回密文
	return cipherText
}

//RSA解密
func RSA_Decrypt(cipherText []byte,path string) []byte{
	//打开文件
	file,err:=os.Open(path)
	if err!=nil{
		panic(err)
	}
	defer file.Close()
	//获取文件内容
	info, _ := file.Stat()
	buf:=make([]byte,info.Size())
	file.Read(buf)
	//pem解码
	block, _ := pem.Decode(buf)
	//X509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err!=nil{
		panic(err)
	}
	//对密文进行解密
	plainText,_:=rsa.DecryptPKCS1v15(rand.Reader,privateKey,cipherText)
	//返回明文
	return plainText
}


func main(){
	message := []byte("hello world")
	cipherText:=RSA_Encrypt(message,"public.pem")
	sEnc := base64.StdEncoding.EncodeToString(cipherText)
	fmt.Printf("enc=[%s]\n", sEnc)
	fmt.Println("加密后为：",string(cipherText))
	//解密
	plainText := RSA_Decrypt(cipherText, "private.pem")
	fmt.Println("解密后为：",string(plainText))

}