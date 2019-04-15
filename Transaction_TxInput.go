package main

import (
	"bytes"
	"fmt"
)

type TXInput struct {
	TxID      []byte  //对应output txId
	Index     int    //存储Txoutput的index里面的索引

	//ScriptSiq string //用户名   私钥  signature  签名脚本 解锁脚本 （）
	Signature []byte //数字签名
	PublicKey []byte //公钥，钱包里面
}

//判断当前txInput消费，和指定的address是否一致
func (txInput *TXInput) UnLockWithAddress (pubKeyHash  []byte) bool {
	//return txInput.ScriptSiq == address
	publicKeyHash := PubKeyHash(txInput.PublicKey)
	fmt.Println("传过来的pubkeyHash",pubKeyHash)
	fmt.Println("PubKeyHash(txInput.PublicKey)",publicKeyHash)
	return bytes.Compare(pubKeyHash,publicKeyHash) == 0

}
