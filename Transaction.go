package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Transaction struct {
	TxID    []byte      //交易ID
	Inputs  []*TXInput  //输入
	Outputs []*TXOutput //输出
}

/**
transaction 分2种
1.一种创始区块创建时的transaction
2.转账时产生的Transaction
 */

func NewCoinBaseTransaction(address string) *Transaction {
	txInput := &TXInput{[]byte{}, -1, "genesis data"}
	txOutput := &TXOutput{10, address}
	txcoinbase := &Transaction{[]byte{}, []*TXInput{txInput}, []*TXOutput{txOutput}}
	txcoinbase.setTxID()
	fmt.Println("genesis block 生成..................")
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

//判断当前交易是否是coinbase
func (tx *Transaction) isCoinbaseTx() bool {
	return tx.Inputs[0].Index == -1 && len(tx.Inputs[0].TxID) == 0
}

//创建交易
func NewSimpleTx(from, to string, amout int64, bc *BlockChain, txs []*Transaction) *Transaction {
	var inputs []*TXInput
	var outputs []*TXOutput
	//够用的Inputs
	balance, spendableUTXOs := bc.FindSpendableUTXOs(from, amout, txs)
	for txId, indexArray := range spendableUTXOs {
		txIdBytes, _ := hex.DecodeString(txId)
		for _, index := range indexArray {
			input := &TXInput{txIdBytes, index, from}
			inputs = append(inputs, input)
		}
	}
	//转账
	output1 := &TXOutput{amout, to}
	outputs = append(outputs, output1)
	//找零
	output2 := &TXOutput{balance - amout, from}
	outputs = append(outputs, output2)

	tx := &Transaction{[]byte{}, inputs, outputs}
	tx.setTxID()
	return tx
}

