package common

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/ripemd160"

	"github.com/mr-tron/base58"
)

type Wallet struct {
	Private *ecdsa.PrivateKey
	PubKey  []byte
}

func NewWallet() *Wallet {
	//建立曲线
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		fmt.Printf("")
	}
	//生成公钥
	pubKeyOrig := privateKey.PublicKey

	pubKey := append(pubKeyOrig.X.Bytes(), pubKeyOrig.Y.Bytes()...)

	return &Wallet{
		Private: privateKey,
		PubKey:  pubKey,
	}
}

func (w *Wallet) NewAddress() string {
	pubKey := w.PubKey
	hash := sha256.Sum256(pubKey)

	rip160hasher := ripemd160.New()
	_, err := rip160hasher.Write(hash[:])

	if err != nil {

	}
	//返回rip160的hash结果
	rip160HashValue := rip160hasher.Sum(nil)

	version := byte(00)
	//拼接version
	payload := append([]byte{version}, rip160HashValue...)

	//checksum
	//两次sha256
	hash1 := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash1[:])
	//校验码
	checkCode := hash2[:4]
	//25字节
	payload = append(payload, checkCode...)
	//go语言有个库，btcd,这个是go语言实现的比特币全节点源码
	s := base58.Encode(payload)
	return s
}
