package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"crypto/ecdsa"
	"encoding/json"
	"crypto/rand"
	"crypto/elliptic"
	"math/big"
)

type Transaction struct {
	TxID    []byte      //交易ID
	Inputs  []*TXInput  //输入
	Outputs []*TXOutput //输出
}

/**
transaction 分2种
1.一种创始区块创建时的transaction
2.转账时产生的Transaction
 */

func NewCoinBaseTransaction(address string) *Transaction {
	txInput := &TXInput{[]byte{}, -1, nil,[]byte{}}
	//txOutput := &TXOutput{10, address}
	txOutput := NewTxOutput(10,address)
	txcoinbase := &Transaction{[]byte{}, []*TXInput{txInput}, []*TXOutput{txOutput}}
	txcoinbase.setTxID()
	fmt.Println("coin base  block 生成..................")
	return txcoinbase
}


func (tx *Transaction) setTxID() {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	buffBytes := bytes.Join([][]byte{IntToHex(time.Now().Unix()), buff.Bytes()}, []byte{})
	hash := sha256.Sum256(buffBytes)
	tx.TxID = hash[:]
}

//判断当前交易是否是coinbase
func (tx *Transaction) isCoinbaseTx() bool {
	return tx.Inputs[0].Index == -1 && len(tx.Inputs[0].TxID) == 0
}

//创建交易
func NewSimpleTx(from, to string, amount int64, bc *BlockChain, txs []*Transaction) *Transaction {
	var inputs []*TXInput
	var outputs []*TXOutput
	//获取钱包
	wallets := NewWallets()
	wallet := wallets.WalletsMap[from]

	//够用的Inputs
	balance, spendableUTXOs := bc.FindSpendableUTXOs(from, amount, txs)
	for txId, indexArray := range spendableUTXOs {
		txIdBytes, _ := hex.DecodeString(txId)
		for _, index := range indexArray {
			input := &TXInput{txIdBytes, index, nil,wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}
	//转账
	//output1 := &TXOutput{amount, to}
	output1 := NewTxOutput(amount,to)
	outputs = append(outputs, output1)
	//找零
	//output2 := &TXOutput{balance - amount, from}

	output2 := NewTxOutput(balance - amount, from)
	fmt.Println("balance",balance)
	fmt.Println("amount",amount)
	fmt.Println("balance - amount",balance - amount)
	fmt.Println("找零",output2.Value)
	outputs = append(outputs, output2)

	tx := &Transaction{[]byte{}, inputs, outputs}
	tx.setTxID()

/*	for _, tx := range txs {
		for _, input := range tx.Inputs {

			fmt.Println("--id",input.TxID)
			fmt.Println("--index",input.Index)
			fmt.Println("--publickey",input.PublicKey)
			fmt.Println("--Signature",input.Signature)
		}
	}
*/
	//创建签名
	bc.SignTransaction(tx,wallet.PrivateKey,txs)
	return tx
}

//签名
//正如上面提到的，为了对一笔交易进行签名，我们需要获取交易输入所引用的输出，因为我们需要存储这些输出的交易。
func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey, preTxs map[string]*Transaction) {
	//1.如果是coinbase交易，无需签名
	if tx.isCoinbaseTx() {
		return
	}
	//2.input没有对应的transaction,无需签名
	for _, input := range tx.Inputs {
		if preTxs[hex.EncodeToString(input.TxID)].TxID == nil {
			log.Panic("当前的input没有对应的transcation")
		}
	}

	//3.获取Transaction的部分数据的副本
	txCopy := tx.TrimmedCopy()
	//4.
	for index, input := range txCopy.Inputs {
		preTx := preTxs[hex.EncodeToString(input.TxID)]
		//为txcopy设置新的交易ID: txID->[]byte{},index,sign-->nil,publicKey--->对应输出的公钥哈希

		input.Signature = nil                                   //
		input.PublicKey = preTx.Outputs[input.Index].PubKeyHash //设置input的公钥为对应output的公钥哈希
		data := txCopy.getData()                                //设置新的txId   []byte{}  --> sha256

		input.PublicKey = nil //

		//签名
		/*
		通过 privKey 对 txCopy.ID 进行签名。
		一个 ECDSA 签名就是一对数字，我们对这对数字连接起来，并存储在输入的 Signature 字段。
		 */
		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, data)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Inputs[index].Signature = signature

	}
}

//获取签名所需要的Transaction的副本
//创建tx的副本：需要剪裁数据
/*
TxID，
[]*TxInput,
	TxInput中，去除sign，publicKey
[]*TxOutput

这个副本包含了所有的输入和输出，但是 TXInput.Signature 和 TXIput.PubKey 被设置为 nil。
 */
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []*TXInput
	var outputs []*TXOutput
	for _, input := range tx.Inputs {
		inputs = append(inputs, &TXInput{input.TxID, input.Index, nil, nil})
	}
	for _, output := range tx.Outputs {
		outputs = append(outputs, &TXOutput{output.Value, output.PubKeyHash})
	}
	txCopy := Transaction{tx.TxID, inputs, outputs}
	return txCopy
}
func (tx *Transaction) getData() []byte {
	txCopy := tx
	txCopy.TxID = []byte{}
	hash := sha256.Sum256(txCopy.Serialize())
	return hash[:]

}
func (tx *Transaction) Serialize() []byte {
	jsonByte, err := json.Marshal(tx)
	if err != nil {
		log.Panic(err)
	}
	return jsonByte
}

//验证数字签名
func (tx *Transaction) Verify(preTxs map[string]*Transaction) bool {
	fmt.Println("开始验证签名。。")
	if tx.isCoinbaseTx() {
		fmt.Println("conibase....")
		return true
	}
	for _, input := range tx.Inputs {
		if preTxs[hex.EncodeToString(input.TxID)].TxID == nil {
			log.Panic("当前的input没有对应的transaction,无法验证。。。")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()
	for index, input := range tx.Inputs {
		fmt.Println("id",input.TxID)
		fmt.Println("index",input.Index)
		fmt.Println("publickey",input.PublicKey)
		fmt.Println("Signature",input.Signature)
		preTx := preTxs[hex.EncodeToString(input.TxID)]
		txCopy.Inputs[index].Signature = nil
		txCopy.Inputs[index].PublicKey = preTx.Outputs[input.Index].PubKeyHash
		data := txCopy.getData()
		txCopy.Inputs[index].PublicKey = nil

		//签名种的s和r
		r := big.Int{}
		s := big.Int{}
		fmt.Println("input.Signature",input.Signature)
		sigLen := len(input.Signature)
		r.SetBytes(input.Signature[:sigLen/2])
		s.SetBytes(input.Signature[sigLen/2:])
		fmt.Println("r-->",r)
		fmt.Println("s-->",s)
		//通过公钥，产生新的s和r ,对比   ？？
		x := big.Int{}
		y := big.Int{}
		keyLen := len(input.PublicKey)
		x.SetBytes(input.PublicKey[:keyLen/2])
		y.SetBytes(input.PublicKey[keyLen/2:])
		fmt.Println("x------>",x)
		fmt.Println("y------>",y)
		//根据椭圆曲线,一技x,y 获取公钥
		//我们使用从input.PUblicKey创建一个 ecdsa.publickey
		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		/**
		这里我们解包存储在txinput.Signature  和Txinput.Pubkey 中的值
		因为一个签名就是一对数字，一个公钥就是一对坐标
		我们之前为了存储将他们连接在一起，现在我们需要对他们进行解包在  crypto/ecdsa 函数中使用。

		验证
		在这里：我们使用从输入提取的公钥创建一个ecdsa.PublicKey,通过传入输入中提取的签名执行了ecdsa.Verity.
		如果所有的输入都背验证，返回true;如果有任何一个失败，返回false
		 */
		 fmt.Println("&rawPubKey",&rawPubKey)
		 fmt.Println("data",data)
		 fmt.Println("&r",&r)
		 fmt.Println("&s",&s)
		if ecdsa.Verify(&rawPubKey, data, &r, &s) == false {
			fmt.Println("验证 不相等。。。。")
			return false
		} 
	}
	return true


}
