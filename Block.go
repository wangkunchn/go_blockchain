package main

import (
	"time"
)

type Block struct {
	Height int64			//区块链高度
	PreBlockHash []byte		//上一个区块的哈希值
	Data [] byte			//交易数据  后期为transaction
	TimeStamp int64			//时间戳
	Hash []byte				//哈希值 32个字节，64个16进制数
	Nonce int64
}

func CreateGenesisBlock(data string) *Block {
	return NewBlock(data, []byte{}, 0)
}

func NewBlock(data string, preBlockHash []byte, height int64) *Block {
	block := &Block{height,preBlockHash,[]byte(data),time.Now().Unix(), nil,0}
	//block.setHash()
	pow := NewProofOfWork(block)
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return block
}