package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"lib/base58"
	"log"
)

//创建一个结构WalletKeyPair密钥对，保存私钥和公钥
//给这个结构提供一个方法GetAddress：私钥->公钥->地址
type WalletKeyPair struct {
	PrivateKey *ecdsa.PrivateKey
	//type PublicKey struct {
	//	elliptic.Curve
	//	X, Y *big.Int
	//}
	//我们可以将公钥的X,Y进行字节流拼接后传输，这样在对端再进行切割还原，好处是方便以后的编码
	PublicKey []byte
}

func NewWalletKeyPair() *WalletKeyPair {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic()
	}
	publicKeyRaw := privateKey.PublicKey
	publicKey := append(publicKeyRaw.X.Bytes(), publicKeyRaw.Y.Bytes()...)
	return &WalletKeyPair{privateKey, publicKey}

}
func (w *WalletKeyPair) GetAdress() string {
	//hash := sha256.Sum256(w.PublicKey)
	////创建一个hash160对象
	////向hash160中write对象
	////做哈希运算
	//rip160Hasher := ripemd160.New()
	//_, err := rip160Hasher.Write(hash[:])
	//if err != nil {
	//	log.Panic()
	//}
	////Sum函数会把我们的结果与Sum参数append到一起，然后返回，我们传入nil,防止数据污染
	////160位，20字节
	//publicHash := rip160Hasher.Sum(nil)
	publicHash := HashPubKey(w.PublicKey)
	version := 0x00
	//21字节的数据
	payload := append([]byte{byte(version)}, publicHash...)
	//first := sha256.Sum256(payload)
	//second := sha256.Sum256(first[:])
	////4字节校验码
	//checksum := second[0:4]
	checksum := CheckSum(payload)
	//25字节
	payload = append(payload, checksum...)
	address := base58.Encode(payload)
	return address

}

func HashPubKey(pubKey []byte) []byte {
	hash := sha256.Sum256(pubKey)
	//创建一个hash160对象
	//向hash160中write对象
	//做哈希运算
	rip160Hasher := ripemd160.New()
	_, err := rip160Hasher.Write(hash[:])
	if err != nil {
		log.Panic()
	}
	//Sum函数会把我们的结果与Sum参数append到一起，然后返回，我们传入nil,防止数据污染
	//160位，20字节
	publicHash := rip160Hasher.Sum(nil)
	return publicHash
}
func CheckSum(payload []byte) []byte {
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])
	//4字节校验码
	checksum := second[0:4]
	return checksum
}
