package main

import (
	"bytes"
	"fmt"
	"time"
)

//cli只用于解析命令
//commands用于实现cli的具体命令

func (cli *CLI) AddBlock(txs []*Transaction) {
	cli.bc.addBlock(txs)
	fmt.Println("添加区块成功！\n")
}

func (cli *CLI) printBlock() {
	it := cli.bc.NewIterator()
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
		//fmt.Printf("Transactions:%v\n", block.Transactions)//TODO
		pow := NewProofOfWork(block)
		fmt.Printf("IsVaild:%v\n", pow.IsValid())
		if bytes.Equal(block.PrevBlockHash, []byte{}) {
			fmt.Println("区块链遍历结束！")
			break
		}
	}

}
