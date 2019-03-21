package main

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"time"
)

type BlockChain struct {
	//Blocks []*Block
	DB  *bolt.DB //数据库对象
	Tip []byte   //最后一个block的hash
}

//创建一个区块链，包含创世区块
func CreateBlockChainWithGenesisBlock(data string)  {
	if dbExists() {
		fmt.Println("数据库已经存在。。")
		return
	}

	genesisBlock := CreateGenesisBlock(data)
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
			fmt.Println("存入创世--pre",genesisBlock.Height,genesisBlock.PreBlockHash)
			fmt.Println("存入创世--hash",genesisBlock.Height,genesisBlock.Hash)
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
func (bc *BlockChain) AddBlockToBlockChain(data string) {
	err := bc.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCK_TABLE_NAME))
		if bucket != nil {
			lastblockHash := bucket.Get(bc.Tip)
			lastBlock := DeserializeBlock(lastblockHash)
			newBlock := NewBlock(data, lastBlock.Hash, lastBlock.Height+1)
			e := bucket.Put(newBlock.Hash, newBlock.Serilalize())
			fmt.Println("存入pre",newBlock.Height,newBlock.PreBlockHash)
			fmt.Println("存入current----",newBlock.Height,newBlock.Hash)
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
		fmt.Printf("\t数据：%s\n", block.Data)
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
			blockchain = &BlockChain{db,hash}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return blockchain
}