package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

//交易输入
//指明交易发起人可支付资金的来源，包含：
//1、引用UTXO所在交易的ID
//2、所消费UTXO在output中的索引
//3、解锁脚本
type TXInput struct {
	//引用utxo所在交易ID（知道在那个房间）
	TXID []byte
	//所消费utxo在output的索引值（具体位置）
	Index int64
	//解锁脚本（签名，公钥）
	Address string
}

//交易输出
//包含资金接收方的相关信息，包含：
//接收金额（数字）
//锁定脚本（对方公钥的哈希，这个哈希可以通过地址反推出来，所以转账知道地址即可）
type TXoutput struct {
	Value   float64
	Address string
}

type Transaction struct {
	TXid []byte //交易ID
	//所有的inputs
	TXInputs []TXInput
	//所有的outputs
	TXoutputs []TXoutput
}

//交易ID
//一般是交易结构的哈希值（参考block序列化）
func (tx *Transaction) SetTXID() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(buffer.Bytes())
	tx.TXid = hash[:]
}

//实现挖矿交易
//特点：只有输出，没有有效的输入（不需要TXInput,包括ID，索引，签名等）
//把挖矿的人传递进来，因为有奖励
func NewCoinbaseTX(miner string, data string) *Transaction {
	inputs := []TXInput{{nil, -1, data}}
	outputs := []TXoutput{{12.5, miner}}
	tx := Transaction{nil, inputs, outputs}
	tx.SetTXID()
	return &tx
}
func (tx *Transaction) IsCoinbase() bool {
	//特点：1、只有1个input 2、引用的id是nil 3、引用的索引是-1
	inputs := tx.TXInputs
	if len(inputs) == 1 && inputs[0].TXID == nil && inputs[0].Index == -1 {
		return true
	}
	return false

}

//参数：
//1、付款人
//2、收款人
//3、转账金额
//4、bc
//创建普通交易的内部逻辑
//1、遍历账本，找到属于付款人的合适的金额，把这个outputs找到
//2、如果找到的outputs的钱不足以转账，创建交易失败
//3、将outputs转成inputs
//4、创建输出，创建一个属于收款人的output
//5、如果有零钱，创建属于付款人的output
//6、设置交易ID
//7、返回交易结构
func NewTransaction(from string, to string, amount float64, bc *Blockchain) *Transaction {
	//map[string][]int64
	//1、遍历账本，找到属于付款人的合适的金额，把这个outputs找到
	validUTXOs := make(map[string][]int64) //标识有用的utxo
	var resValue float64                   //这些utxo存储的金额
	//第一部分，找到所需要的UTXO的集合
	validUTXOs /*本次支付所需要的utxo集合*/, resValue /*返回utxos所包含的金额*/ = bc.FindSuitableUTXOs(from, amount)
	//2、如果找到的outputs的钱不足以转账，创建交易失败
	if resValue < amount {
		fmt.Printf("余额不足，交易失败！！！\n")
		os.Exit(1)
	}
	////3、将outputs转成inputs
	var inputs []TXInput
	var outputs []TXoutput
	//第二部分，input的创建
	for txid, indexes /*0x333*/ := range validUTXOs {
		for _, i /*索引：0,1*/ := range indexes {
			input := TXInput{[]byte(txid), i, from}
			inputs = append(inputs, input)
		}
	}
	//4、创建输出，创建一个属于收款人的output
	output := TXoutput{amount, to}
	outputs = append(outputs, output)
	//5、如果有零钱，创建属于付款人的output
	if resValue > amount {
		output1 := TXoutput{resValue - amount, from}
		outputs = append(outputs, output1)
	}

	//创建交易
	tx := Transaction{nil, inputs, outputs}
	//6、设置交易ID
	tx.SetTXID()
	//7、返回交易结构
	return &tx

}
