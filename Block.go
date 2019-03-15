package main

import (
	"time"
	"bytes"
	"crypto/sha256"
)

type Block struct {
	Height int64			//区块链高度
	PreBlockHash []byte		//上一个区块的哈希值
	Data [] byte			//交易数据  后期为transaction
	TimeStamp int64			//时间戳
	Hash []byte				//哈希值 32个字节，64个16进制数
}

func (block *Block) setHash() {
	//每个数据都抓成[]byte
	heightBytes := IntToHex(block.Height)
	timeBytes := IntToHex(block.TimeStamp)
	//拼接
    s := bytes.Join([][]byte{heightBytes, block.PreBlockHash, block.Data, timeBytes}, []byte{})
	//hash
	sum256 := sha256.Sum256(s)
	block.Hash = sum256[:]
}

func CreateGenesisBlock(data string) *Block {
	return NewBlock(data, []byte{}, 0)
}

func NewBlock(data string, preBlockHash []byte, height int64) *Block {
	block := &Block{height,preBlockHash,[]byte(data),time.Now().Unix(), nil}
	block.setHash()
	return block
}