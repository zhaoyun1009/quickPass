package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
)

func GenerateKey() (publicKey, privateKey string, err error) {
	private, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}
	privateKey, err = dumpPrivateKeyBase64(private)
	if err != nil {
		return "", "", err
	}

	publicKey, err = dumpPublicKeyBase64(&private.PublicKey)
	if err != nil {
		fmt.Println(err)
	}
	return privateKey, publicKey, nil
}

func RSASign(data []byte, privateKey string) (string, error) {
	// 1、选择hash算法，对需要签名的数据进行hash运算
	hash := sha256.New()
	hash.Write(data)
	// 2、解析出私钥对象
	rsaPrivateKey, err := loadPrivateKeyBase64(privateKey)
	if err != nil {
		return "", err
	}
	// 3、RSA数字签名（参数是随机数、私钥对象、哈希类型、签名文件的哈希串，生成bash64编码）
	bytes, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, crypto.SHA256, hash.Sum(nil))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func RSAVerify(data []byte, base64Sign, publicKey string) error {
	// 1、对base64编码的签名内容进行解码，返回签名字节
	sign, err := base64.StdEncoding.DecodeString(base64Sign)
	if err != nil {
		return err
	}
	// 2、选择hash算法，对需要签名的数据进行hash运算
	hash := sha256.New()
	hash.Write(data)
	// 3、读取公钥文件，解析出公钥对象
	rsaPublicKey, err := loadPublicKeyBase64(publicKey)
	if err != nil {
		return err
	}
	// 4、RSA验证数字签名（参数是公钥对象、哈希类型、签名文件的哈希串、签名后的字节）
	return rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, hash.Sum(nil), sign)
}

func dumpPrivateKeyBase64(privatekey *rsa.PrivateKey) (string, error) {
	keybytes := x509.MarshalPKCS1PrivateKey(privatekey)
	keybase64 := base64.StdEncoding.EncodeToString(keybytes)
	return keybase64, nil
}

func dumpPublicKeyBase64(publickey *rsa.PublicKey) (string, error) {
	keybytes, err := x509.MarshalPKIXPublicKey(publickey)
	if err != nil {
		return "", err
	}

	keybase64 := base64.StdEncoding.EncodeToString(keybytes)
	return keybase64, nil
}

// Load private key from base64
func loadPrivateKeyBase64(base64key string) (*rsa.PrivateKey, error) {
	keybytes, err := base64.StdEncoding.DecodeString(base64key)
	if err != nil {
		return nil, fmt.Errorf("base64 decode failed, error=%s\n", err.Error())
	}

	privatekey, err := x509.ParsePKCS1PrivateKey(keybytes)
	if err != nil {
		return nil, errors.New("parse private key error!")
	}

	return privatekey, nil
}

func loadPublicKeyBase64(base64key string) (*rsa.PublicKey, error) {
	keybytes, err := base64.StdEncoding.DecodeString(base64key)
	if err != nil {
		return nil, fmt.Errorf("base64 decode failed, error=%s\n", err.Error())
	}

	pubkeyinterface, err := x509.ParsePKIXPublicKey(keybytes)
	if err != nil {
		return nil, err
	}

	publickey := pubkeyinterface.(*rsa.PublicKey)
	return publickey, nil
}
