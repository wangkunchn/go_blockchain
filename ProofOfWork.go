package main

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

const TargetBit = 16 //20 24

type ProofOfWork struct {
	Block *Block //验证的区块
	Target *big.Int //大整数存储，目标哈希
}


func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, 256 - TargetBit)

	return &ProofOfWork{block, target}
}

//挖矿 生成hash,nonce
func (pow *ProofOfWork) Run() ([]byte, int64) {
 for {
 	var nonce int64 = 0
 	hashInt := new(big.Int)
 	var blockHash []byte
 	for {
		blockHash = pow.prepareData(nonce)
		fmt.Printf("\r%d: %x\n",nonce,blockHash)
		hashInt.SetBytes(blockHash)

		if pow.Target.Cmp(hashInt) == 1{
			break
		}
		nonce++
	}
	return blockHash,nonce
 }
}

//hash 数据
func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	data := bytes.Join([][]byte{
		IntToHex(pow.Block.Height),
		pow.Block.Data,
		pow.Block.PreBlockHash,
		IntToHex(pow.Block.TimeStamp),
		IntToHex(nonce),
	},[]byte{})
	blockHash := sha256.Sum256(data)
	return blockHash[:]
}
//验证
func (pow *ProofOfWork) isValid() bool {
	hashInt := new(big.Int)
	hashInt.SetBytes(pow.Block.Hash)
	return pow.Target.Cmp(hashInt) == 1
}

/**
 A：Bits = "0x17502ab7"

 B：exponent指数，exponent = 0x17

 C：coefficient系数，coefficient = 0x502ab7

 D：target = coefficient * Math.Pow(2, 8 * (exponent - 3))

 E：目标hash：000000000000000000502ab700000000d6420b16625d309c4561290000000000

 F：实际hash：00000000000000000041ff1cfc5f15f929c1a45d262f88e4db83680d90658c0c
 */