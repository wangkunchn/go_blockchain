package main

import (
	"flag"
	"os"
	"log"
	"fmt"
)

type CLI struct {
	//BLockchain *blockchain
}

func (cli *CLI) Run() {
	isVaild()
	//1.创建flagset标签
	createBlockChainCmd := flag.NewFlagSet("creatBlockChain", flag.ExitOnError)
	addBlockCmd := flag.NewFlagSet("addBlock", flag.ExitOnError)
	printChainCmd:= flag.NewFlagSet("printChain", flag.ExitOnError)
	
	//2.设置标签后面的参数
	createBlockChainData := createBlockChainCmd.String("address", "Genesis block data..", "创始区块交易数据")
	addBlockData := addBlockCmd.String("data", "helloworld..", "交易数据")

	//3.解析
	switch os.Args[1] {
	case "creatBlockChain" :
		err := createBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "addBlock" :
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printChain" :
		err :=printChainCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)//推出
	}

	if createBlockChainCmd.Parsed() {
		if *createBlockChainData == "" {
			printUsage()
			os.Exit(1)
		}
		fmt.Println(*createBlockChainData)
		cli.createGenesisBlock(*createBlockChainData)
	}
	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			printUsage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}
	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func (cli *CLI) createGenesisBlock(address string) {
	CreateBlockChainWithGenesisBlock(address)
}
func (cli *CLI) addBlock(data string) {
	bc := GetBlockChainObject()
	if bc == nil {
		fmt.Println("没有创始区块可以添加。。")
		return
	}
	defer bc.DB.Close()
	//bc.AddBlockToBlockChain(data)
}

func (cli *CLI) printChain() {
	bc := GetBlockChainObject()
	if bc == nil {
		fmt.Println("没有区块链可以打印。。")
		return
	}
	defer bc.DB.Close()
	bc.PringChains()

}

func isVaild() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}


func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tcreatBlockChain -address DATA -- 创建coinbase")
	fmt.Println("\taddBlock -data Data -- 交易数据")
	fmt.Println("\tprintChain -- 输出信息")
}
