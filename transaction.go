package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
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
func NewCoinbaseTX(miner string) *Transaction {
	//TODO
	inputs := []TXInput{{nil, 1, "hello"}}
	outputs := []TXoutput{{12.5, "miner"}}
	tx := Transaction{nil, inputs, outputs}
	tx.SetTXID()
	return &tx

}
