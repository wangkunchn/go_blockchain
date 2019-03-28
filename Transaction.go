package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"github.com/boltdb/bolt"
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

//挖矿 产生新区块
func (bc *BlockChain) MineNewBlock(from, to, amout []string) {
	/**
			go_blockchain send -from '["wangkun"]' -to '["baby"]' -amout '["4"]'

			["wangkun"]		["baby"]		["4"]

			1.new tx/txs
			2.new block
			3.block to blockChain   block加入数据库
	 */
	var txs []*Transaction
	for i := 0; i < len(from); i++ {
		amountInt, _ := strconv.ParseInt(amout[i], 10, 64)
		tx := NewSimpleTx(from[i], to[i], amountInt, bc, txs)
		txs = append(txs, tx)
	}
	var preBlock *Block
	var newBlock *Block
	bc.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCK_TABLE_NAME))
		if bucket != nil {
			lastHash := bucket.Get([]byte(LAST_BLOCK_HASH))
			preBlockBytes := bucket.Get(lastHash)
			preBlock = DeserializeBlock(preBlockBytes)
		}
		return nil
	})

	newBlock = NewBlock(txs, preBlock.Hash, preBlock.Height+1)

	//打开数据库 添加block
	bc.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCK_TABLE_NAME))
		if bucket != nil {
			newBlockBytes := newBlock.Serilalize()
			bucket.Put(newBlock.Hash,newBlockBytes)
			bucket.Put([]byte(LAST_BLOCK_HASH),newBlock.Hash)
			bc.Tip = newBlock.Hash
			fmt.Println(newBlock.Height," : new block Hash-->",newBlock.Hash)
		}
		return nil
	})
}
