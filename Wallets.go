package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"log"
	"encoding/gob"
	"crypto/elliptic"
	"bytes"
)

type Wallets struct {
	WalletsMap map[string]*Wallet
}

const walletFile = "Wallets.dat"

func NewWallets() *Wallets {
	//1.判断文件是否存在
	_, err := os.Stat(walletFile);
	if os.IsNotExist(err) {
		fmt.Println("文件不存在")
		wallets := &Wallets{}
		wallets.WalletsMap = make(map[string]*Wallet)
		return wallets
	}
	//2.否则读取文件种的数据
	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}
	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}
	return &wallets
}

func (ws *Wallets) CreateNewWallet()  {
	wallet := NewWallet()
	fmt.Printf("创建钱包地址：%s\n", wallet.GetAddress())
	ws.WalletsMap[string(wallet.GetAddress())] = wallet
	ws.SaveWallets()
}
/**
持久化存储
gob是Golang包自带的一个数据结构序列化的编码/解码工具
 */
func (ws *Wallets) SaveWallets() {
	var content bytes.Buffer
	//注册的目的，为了可以序列化任何类型，wallet结构体种有接口类型。将接口进行注册
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}
	//将序列化后的数据content 写入文件，原来内容会被覆盖
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}