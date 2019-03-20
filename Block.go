package main

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
	"io"
	"fmt"
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
	return NewBlock(data, make([] byte,32,32), 0)
}

func NewBlock(data string, preBlockHash []byte, height int64) *Block {
	block := &Block{height,preBlockHash,[]byte(data),time.Now().Unix(), nil,0}
	pow := NewProofOfWork(block)
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	fmt.Println("newBlock",block.Height)
	return block
}

func (block *Block) Serilalize() []byte  {
	//1.创建一个buffer
	var result bytes.Buffer
	//2.创建一个编码器
	encoder := gob.NewEncoder(&result)
	//3.编码-->打包
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func DeserializeBlock(blockBytes []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err == io.EOF {
		return &block
	} else if err != nil {
		log.Panic(err)
	}
	return &block

}