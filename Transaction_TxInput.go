package main

type TXInput struct {
	TxID      []byte  //对应output txId
	Index     int    //存储Txoutput的index里面的索引
	ScriptSiq string //用户名   私钥  signature  签名脚本 解锁脚本 （）
}

//判断当前txInput消费，和指定的address是否一致
func (txInput *TXInput) UnLockWithAddress (address string) bool {
	return txInput.ScriptSiq == address

}
