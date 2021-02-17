package main

import (
	"bytes"
	"fmt"
	"time"
)

//cli只用于解析命令
//commands用于实现cli的具体命令

func (cli *CLI) CreateBlockchain(address string) {
	if !IsValidAddress(address) {
		fmt.Printf("%s是无效地址！\n", address)
		return
	}
	bc := CreateBlockchain(address)
	if bc == nil {
		return
	}
	defer bc.db.Close()
	fmt.Println("创建区块链成功！")
}

func (cli *CLI) AddBlock(txs []*Transaction) {
	bc := NewBlockchain()
	if bc == nil {
		return
	}
	defer bc.db.Close()
	bc.AddBlock(txs)
	fmt.Printf("添加区块成功！\n")
}

func (cli *CLI) printBlock() {
	bc := NewBlockchain()
	if bc == nil {
		return
	}
	defer bc.db.Close()
	it := bc.NewIterator()
	for {
		block := it.Next()
		fmt.Println("-------------------------------")
		fmt.Printf("Version:%x\n", block.Version)
		fmt.Printf("PrevBlockHash:%x\n", block.PrevBlockHash)
		fmt.Printf("MerkleRoot:%x\n", block.MerkleRoot)
		//优化时间打印，定义时间显示格式
		timeFormat := time.Unix(int64(block.Timestamp), 0).Format("2006-01-02 15:01:05")
		fmt.Printf("Timestamp:%s\n", timeFormat)
		fmt.Printf("Difficulity:%x\n", block.Difficulity)
		fmt.Printf("Nonce:%x\n", block.Nonce)
		fmt.Printf("Hash:%x\n", block.Hash)
		fmt.Printf("Data:%v\n", block.Transactions[0].TXInputs[0].PublicKey)
		pow := NewProofOfWork(block)
		fmt.Printf("IsVaild:%v\n", pow.IsValid())
		if bytes.Equal(block.PrevBlockHash, []byte{}) {
			fmt.Println("区块链遍历结束！")
			break
		}
	}

}

func (cli *CLI) Send(from, to string, amount float64, miner string, data string) {
	if !IsValidAddress(from) {
		fmt.Printf("%s是无效地址！\n", from)
		return
	}
	if !IsValidAddress(to) {
		fmt.Printf("%s是无效地址！\n", to)
		return
	}
	if !IsValidAddress(miner) {
		fmt.Printf("%s是无效地址！\n", miner)
		return
	}
	bc := NewBlockchain()
	if bc == nil {
		return
	}
	defer bc.db.Close()
	//创建挖矿交易
	coinbase := NewCoinbaseTX(miner, data)
	//创建普通交易
	txs := []*Transaction{coinbase}
	tx := NewTransaction(from, to, amount, bc)

	if tx != nil {
		txs = append(txs, tx)
	} else {
		fmt.Printf("发现无效交易，过滤！\n")
	}
	//添加到区块
	bc.AddBlock([]*Transaction{coinbase, tx})

	fmt.Printf("挖矿成功！\n")
}
func (cli *CLI) GetBalance(address string) {
	if !IsValidAddress(address) {
		fmt.Printf("%s是无效地址！\n", address)
		return
	}
	bc := NewBlockchain()
	if bc == nil {
		return
	}
	defer bc.db.Close()
	bc.GetBalance(address)
}
func (cli *CLI) CreateWallet() {
	w := NewWallets()
	address := w.CreateWallet()
	fmt.Printf("新的钱包地址为：%s\n", address)
}
func (cli *CLI) ListAddresses() {
	ws := NewWallets()
	addresses := ws.ListAddress()
	for _, address := range addresses {
		fmt.Printf("address : %s \n", address)
	}
}
func (cli *CLI) PrintTx() {
	bc := NewBlockchain()
	if bc == nil {
		return
	}
	defer bc.db.Close()
	it := bc.NewIterator()
	for {
		block := it.Next()
		fmt.Printf("\n+++++++++++++++++++++新的区块+++++++++++++++++++++\n")
		for _, tx := range block.Transactions {
			fmt.Printf("tx : %v\n", tx)

		}
		if bytes.Equal(block.PrevBlockHash, []byte{}) {
			fmt.Println("区块链遍历结束！")
			break
		}
	}

}
