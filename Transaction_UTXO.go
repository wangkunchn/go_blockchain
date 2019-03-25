package main
//unspent transaction outputss 未交易输出
type UTXO struct {
	TxID []byte	//当前transaction 的 ID
	Index int	//索引
	Output *TXOutput//要使用的output
}

