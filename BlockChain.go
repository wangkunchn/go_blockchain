package main

type BlockChain struct {
	Blocks []*Block
}

//创建一个区块链，包含创世区块
func CreateBlockChainWithGenesisBlock(data string) *BlockChain {
	genesisBlock := CreateGenesisBlock(data)
	return &BlockChain{[]*Block{genesisBlock}}
}

//添加区块到链中
func (bc *BlockChain)AddBlockToBlockChain(newBlock *Block)  {
	bc.Blocks = append(bc.Blocks, newBlock)
}