package main

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"time"
	"encoding/hex"
	"os"
	"strconv"
)

type BlockChain struct {
	//Blocks []*Block
	DB  *bolt.DB //数据库对象
	Tip []byte   //最后一个block的hash
}

//创建一个区块链，包含创世区块
func CreateBlockChainWithGenesisBlock(address string) {
	if dbExists() {
		fmt.Println("数据库已经存在。。")
		return
	}

	txcoinbase := NewCoinBaseTransaction(address)
	genesisBlock := CreateGenesisBlock([]*Transaction{txcoinbase})
	db, err := bolt.Open(DB_NAME, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	//defer db.close()
	err2 := db.Update(func(tx *bolt.Tx) error {
		bucket, e := tx.CreateBucketIfNotExists([]byte(BLOCK_TABLE_NAME))
		if e != nil {
			log.Panic(e)
		}
		if bucket != nil {
			err := bucket.Put(genesisBlock.Hash, genesisBlock.Serilalize())
			fmt.Println("存入创世--pre", genesisBlock.Height, genesisBlock.PreBlockHash)
			fmt.Println("存入创世--hash", genesisBlock.Height, genesisBlock.Hash)
			if err != nil {
				log.Panic("存储创始区块有误。。")
			}
			err2 := bucket.Put([]byte(LAST_BLOCK_HASH), genesisBlock.Hash)
			if err2 != nil {
				log.Panic("last block hash 更新有误")
			}
		}
		return nil
	})
	if err2 != nil {
		log.Panic(err2)
	}

}

//添加区块到链中
func (bc *BlockChain) AddBlockToBlockChain(txs []*Transaction) {
	err := bc.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCK_TABLE_NAME))
		if bucket != nil {
			lastblockHash := bucket.Get(bc.Tip)
			lastBlock := DeserializeBlock(lastblockHash)
			newBlock := NewBlock(txs, lastBlock.Hash, lastBlock.Height+1)
			e := bucket.Put(newBlock.Hash, newBlock.Serilalize())
			fmt.Println("存入pre", newBlock.Height, newBlock.PreBlockHash)
			fmt.Println("存入current----", newBlock.Height, newBlock.Hash)
			if e != nil {
				log.Panic("add block to blockchain err!!!!。。..")
			}
			e2 := bucket.Put([]byte(LAST_BLOCK_HASH), newBlock.Hash)
			if e2 != nil {
				log.Panic("last block hash 更新有误。。。。。")
			}
			bc.Tip = newBlock.Hash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{bc.Tip, bc.DB}
}

func (bc *BlockChain) PringChains() {
	iterator := bc.Iterator()
	var count = 0
	for {
		block := iterator.Next()
		count++
		fmt.Println(count)
		fmt.Printf("第%d个区块的信息：\n", count)
		//获取当前hash对应的数据，并进行反序列化
		fmt.Printf("\t高度：%d\n", block.Height)
		fmt.Printf("\t上一个区块的hash：%x\n", block.PreBlockHash)
		fmt.Printf("\t当前的hash：%x\n", block.Hash)
		fmt.Println("\t交易：")
		for _, tx := range block.Txs {
			fmt.Printf("\t\t交易ID：%x\n", tx.TxID)
			fmt.Println("\t\tVins:")
			for _, in := range tx.Inputs {
				fmt.Printf("\t\t\tTxID:%x\n", in.TxID)
				fmt.Printf("\t\t\tVout:%d\n", in.Index)
				fmt.Printf("\t\t\tScriptSiq:%s\n", in.ScriptSiq)
			}
			fmt.Println("\t\tVouts:")
			for _, out := range tx.Outputs {
				fmt.Printf("\t\t\tvalue:%d\n", out.Value)
				fmt.Printf("\t\t\tScriptPubKey:%s\n", out.ScriptPubKey)
			}
		}
		//fmt.Printf("\t时间：%v\n", block.TimeStamp)
		fmt.Printf("\t时间：%s\n", time.Unix(block.TimeStamp, 0).Format("2006-01-02 15:04:05"))
		fmt.Printf("\t次数：%d\n", block.Nonce)

		//知道block height 0
		if block.Height == 0 {
			break
		}
	}

}

func GetBlockChainObject() *BlockChain {
	if !dbExists() {
		fmt.Println("数据库不存在无法返回 blockchain")
		return nil
	}
	db, e := bolt.Open(DB_NAME, 0600, nil)
	if e != nil {
		log.Panic(e)
	}
	//defer db.close()
	var blockchain *BlockChain
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCK_TABLE_NAME))
		if bucket != nil {
			hash := bucket.Get([]byte(LAST_BLOCK_HASH))
			blockchain = &BlockChain{db, hash}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return blockchain
}

//所有该address的 UTXO
func (bc *BlockChain) AllUTXOs(address string, txs []*Transaction) []*UTXO {
	var allUTXOs []*UTXO
	spentTXOs := make(map[string][]int)

	//1.添加先从txs遍历，查找未花费
	for i := len(txs) - 1; i >= 0; i-- {
		allUTXOs = calculateUTXO(address, txs[i], spentTXOs, allUTXOs)
	}
	//2.遍历数据库，获取每个块中的Transaction,找到未花费的Output
	iterator := bc.Iterator()
	for {
		block := iterator.Next()
		for i := len(block.Txs) - 1; i >= 0; i-- {
			allUTXOs = calculateUTXO(address, block.Txs[i], spentTXOs, allUTXOs)
		}
		if block.Height == 0 {
			break
		}
	}

	return allUTXOs
}

func calculateUTXO(address string, tx *Transaction, spentTXOs map[string][]int, utxos []*UTXO) []*UTXO {
	//1.先得到spentUTXOs
	if !tx.isCoinbaseTx() {
		for _, input := range tx.Inputs {
			if input.UnLockWithAddress(address) {
				key := hex.EncodeToString(input.TxID)
				spentTXOs[key] = append(spentTXOs[key], input.Index)
			}
		}
	}
	//2.遍历UTXOs 满足的address,将满足address的 spentTXOs 排除掉
output:
	for i, output := range tx.Outputs {
		if output.UnLockWithAddress(address) {
			if len(spentTXOs) != 0 {

				for key, indexs := range spentTXOs {
					if key == hex.EncodeToString(tx.TxID) {
						for _, index := range indexs {
							if index == i {
								//已花费 对应的  未花费，不加入数组
								continue output
							}
						}
					}
				}
				utxo := &UTXO{tx.TxID, i, output}
				utxos = append(utxos, utxo)
			} else {
				//未花费
				utxo := &UTXO{tx.TxID, i, output}
				utxos = append(utxos, utxo)
			}
		}
	}
	return utxos
}
//转账时 找到部分可用utxo
func (bc *BlockChain) FindSpendableUTXOs(from string,amount int64, txs []*Transaction) (int64, map[string][]int) {
	//获取所有utxo  遍历 返回值： map[hash]{indexs}
	var balance int64
	utxos := bc.AllUTXOs(from, txs)
	spendableUTXOs := make(map[string][]int)
	for _, utxo := range utxos {
		balance += utxo.Output.Value
		idHash := hex.EncodeToString(utxo.TxID)
		spendableUTXOs[idHash] = append(spendableUTXOs[idHash], utxo.Index)
		if balance >= amount {
			//用几个utxo就够了
			break
		}
	}
	fmt.Println(from,":balance:-->", balance)
	if balance < amount {
		fmt.Printf("%s 余额不足。。总额：%d，需要：%d\n", from,balance,amount)
		os.Exit(1)
	}
	return balance,spendableUTXOs
}

func (bc *BlockChain) getBalance(address string,txs []*Transaction) int64 {
	var balance int64
	utxos := bc.AllUTXOs(address, txs)
	for _,utxo := range utxos {
		balance += utxo.Output.Value
	}
	return balance
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
