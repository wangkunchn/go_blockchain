package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
	"crypto/sha256"
)

type Transaction struct {
	TxID []byte		//交易ID
	Vins []*TXInput 		//输入
	Vouts []*TXOutput		//输出
}


/**
transaction 分2种
1.一种创始区块创建时的transaction
2.转账时产生的Transaction
 */

func NewCoinBaseTransaction(address string) *Transaction {
	txInput := &TXInput{[]byte{}, -1, "genesis data"}
	txOutput := &TXOutput{10, address}
	txcoinbase := &Transaction{[]byte{},[]*TXInput{txInput}, []*TXOutput{txOutput}}
	txcoinbase.setTxID()
	return txcoinbase
}

func (tx *Transaction) setTxID() {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	buffBytes := bytes.Join([][]byte{IntToHex(time.Now().Unix()), buff.Bytes()}, []byte{})
	hash := sha256.Sum256(buffBytes)
	tx.TxID = hash[:]
}