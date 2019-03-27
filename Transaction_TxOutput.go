package main

type TXOutput struct {
	Value int64
	ScriptPubKey string //用户名  公钥，锁定脚本 （里面有用户的address）
}

//判断当前txOutput消费，和指定的address是否一致
func (txOutput *TXOutput) UnLockWithAddress(address string) bool {
	return txOutput.ScriptPubKey == address
}