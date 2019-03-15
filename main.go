package main

import (
	"fmt"
)

func main() {
	blockchain := CreateBlockChainWithGenesisBlock("创始区块")
	preBlock0 := blockchain.Blocks[len(blockchain.Blocks)-1]
	blockchain.AddBlockToBlockChain(NewBlock("第2块",preBlock0.Hash,preBlock0.Height + 1))
	preBlock1 := blockchain.Blocks[len(blockchain.Blocks)-1]
	blockchain.AddBlockToBlockChain(NewBlock("第3块",preBlock1.Hash,preBlock1.Height + 1))
	preBlock2 := blockchain.Blocks[len(blockchain.Blocks)-1]
	blockchain.AddBlockToBlockChain(NewBlock("第4块",preBlock2.Hash,preBlock2.Height + 1))
	preBlock3 := blockchain.Blocks[len(blockchain.Blocks)-1]
	blockchain.AddBlockToBlockChain(NewBlock("第5块",preBlock3.Hash,preBlock3.Height + 1))
	fmt.Println(string(blockchain.Blocks[0].Data),blockchain.Blocks[0].Height)
	for _, block := range blockchain.Blocks {
		fmt.Printf("Prev. hash: %x\n", block.PreBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
	}
}
