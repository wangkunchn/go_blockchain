package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"bytes"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey //私钥
	PublicKey  []byte           //公钥
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	/**
	1.通过椭圆曲线算法，随机产生私钥
	2.根据私钥生成公钥
	 */
	curve := elliptic.P256() //椭圆曲线算法，得到一个椭圆曲线值，全称：SECP256k1
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return *privateKey, publicKey
}

func NewWallet() *Wallet {
	privateKey, publicKey := newKeyPair()
	return &Wallet{privateKey, publicKey}
}

const version = byte(0x00)
const addressChecksumLen = 4

func (w *Wallet) GetAddress() []byte {
	//1. publicKey--> sha256 ripemd160
	pubKeyHash := PubKeyHash(w.PublicKey)
	//2. 版本号+ pubkeyhash + 校验码  --> base58
	versioned_payload := append([]byte{version}, pubKeyHash...)
	checkSumBytes := CheckSum(versioned_payload)
	full_payload := append(versioned_payload, checkSumBytes...)
	address := Base58Encode(full_payload)
	return address
}

//获取验证码：将公钥 sha256 两次，取前四位返回
func CheckSum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:addressChecksumLen]
}

//一次sha256 一次ripemd160 得到publicKeyHash
func PubKeyHash(publicKey []byte) []byte {
	hasher := sha256.New()
	hasher.Write(publicKey)
	hash := hasher.Sum(nil)

	ripemder := ripemd160.New()
	ripemder.Write(hash)
	pubkeyHash := ripemder.Sum(nil)

	return pubkeyHash
}

//判断地址是否有效
func IsValidForAddress(address []byte) bool {
	full_payload := Base58Encode(address)
	checkSumBytes := full_payload[len(full_payload)-addressChecksumLen:]
	version_payload := full_payload[:len(full_payload)-addressChecksumLen]
	CheckBytes := CheckSum(version_payload)
	if bytes.Compare(checkSumBytes, CheckBytes) == 0 {
		return true
	}
	return false
}
