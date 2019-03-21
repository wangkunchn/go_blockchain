package main

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockChainIterator struct {
	CurrentHash []byte //当前区块hash
	DB *bolt.DB			//数据库
}

func (i *BlockChainIterator) Next() *Block {
	block := new(Block)
	err := i.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCK_TABLE_NAME))
		if bucket != nil {
			blockHash := bucket.Get(i.CurrentHash)
			block = DeserializeBlock(blockHash)
			//更新iterator 最新hash
			i.CurrentHash = block.PreBlockHash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return block
}
