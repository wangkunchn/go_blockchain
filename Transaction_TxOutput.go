package main

import "bytes"

type TXOutput struct {
	Value int64
	//ScriptPubKey string //用户名  公钥，锁定脚本 （里面有用户的address）
	PubKeyHash []byte //公钥
}

//判断当前txOutput消费，和指定的address是否一致
func (txOutput *TXOutput) UnLockWithAddress(address string) bool {
	//return txOutput.ScriptPubKey == address
	fullPayloadHash := Base58Decode([]byte(address))
	pubKeyHash := fullPayloadHash[1 : len(fullPayloadHash)-4]
	return bytes.Compare(txOutput.PubKeyHash, pubKeyHash) == 0
}

func NewTxOutput(value int64, address string) *TXOutput {
	txOutput := &TXOutput{value, nil}
	txOutput.Lock(address)
	return txOutput

}

func (txOutput *TXOutput) Lock(address string) {
	publicKeyHash := Base58Decode([]byte(address))
	txOutput.PubKeyHash = publicKeyHash[1 : len(publicKeyHash)-4]
}
