package main

type TXInput struct {
	TxID []byte
	Vout int //存储Txoutput的vout里面的索引
	ScriptSiq string //用户名   私钥  signature  签名脚本 解锁脚本 （）
}
