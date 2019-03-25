package main

type TXOutput struct {
	Value int64
	ScriptPubKey string //用户名  公钥，锁定脚本 （里面有用户的address）
}
