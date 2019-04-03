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
	createBlockChainCmd := flag.NewFlagSet("createBlockChain", flag.ExitOnError)
	//addBlockCmd := flag.NewFlagSet("addBlock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)
	sendTxCmd := flag.NewFlagSet("send", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getBalance", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createWallet", flag.ExitOnError)
	addressListCmd := flag.NewFlagSet("addressList", flag.ExitOnError)

	//2.设置标签后面的参数
	createBlockChainData := createBlockChainCmd.String("address", "Genesis block data..", "创始区块交易数据")
	//addBlockData := addBlockCmd.String("data", "helloworld..", "交易数据")
	fromData := sendTxCmd.String("from", "", "转账源地址")
	toData := sendTxCmd.String("to", "", "转账目标地址")
	amountData := sendTxCmd.String("amount", "", "转账金额")
	getBalanceData := getBalanceCmd.String("address", "", "要查询的地址")

	//3.解析
	switch os.Args[1] {
	case "creatBlockChain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createWallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getBalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printChain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendTxCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "addressList":
		err := addressListCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1) //推出
	}

	if createBlockChainCmd.Parsed() {
		if *createBlockChainData == "" {
			printUsage()
			os.Exit(1)
		}
		fmt.Println(*createBlockChainData)
		cli.createGenesisBlock(*createBlockChainData)
	}
	if getBalanceCmd.Parsed() {
		if *getBalanceData == "" {
			fmt.Println("查询地址不能为空")
			printUsage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceData)
	}
	if printChainCmd.Parsed() {
		cli.printChain()
	}
	if sendTxCmd.Parsed() {
		if *fromData == "" || *toData == "" || *amountData == "" {
			printUsage()
			os.Exit(1)
		}
		fmt.Println(*fromData)
		fmt.Println(*toData)
		fmt.Println(*amountData)
		from := JsonToArray(*fromData)
		to := JsonToArray(*toData)
		amount := JsonToArray(*amountData)

		cli.send(from, to, amount)
	}
	if createWalletCmd.Parsed() {
		cli.createWallet()
	}
	if addressListCmd.Parsed() {
		cli.addressList()
	}
}

func (cli *CLI) createGenesisBlock(address string) {
	CreateBlockChainWithGenesisBlock(address)
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
func (cli *CLI) send(from, to, amount []string) {
	if !dbExists() {
		fmt.Println("数据库不存在。。。")
		os.Exit(1)
	}
	bc := GetBlockChainObject()
	bc.MineNewBlock(from, to, amount)
	defer bc.DB.Close()
}
func (cli *CLI) getBalance(address string) {
	bc := GetBlockChainObject()
	if bc == nil {
		fmt.Println("数据库不存在。。")
		return
	}
	defer bc.DB.Close()
	balance := bc.getBalance(address, []*Transaction{})
	fmt.Printf("%s,一共有%d个Token\n",address,balance)
}
func (cli *CLI) createWallet() {
	wallets := NewWallets()
	wallets.CreateNewWallet()
}
func (cli *CLI) addressList() {
	wallets := NewWallets()
	for address := range wallets.WalletsMap {
		fmt.Println("address:",address)

	}
}

func isVaild() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\t createBlockChain -address DATA -- 创建coinbase")
	fmt.Println("\t printChain -- 输出信息")
	fmt.Println("\t createWallet -- 创建钱包地址")
	fmt.Println("\t addressList -- 地址列表")
	fmt.Println("\t send -from Data -to Data -amount Data --转账")
	fmt.Println("\t getBalance -address Data -- 查询地址余额")
}
