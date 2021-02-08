package main

import (
	"fmt"
)

//1、定义结构，区块链的字段比正常少
//1、当前哈希值
//2、前一区块的哈希值
//3、数据
//2、创建区块
//3、生成哈希
//4、引入区块链
//5、添加区块
//6、重构代码

type Block struct {
	PrevBlockHash []byte //前一哈希值
	Hash          []byte //当前哈希值
	Data          []byte //数据，目前使用字节流，v4开始使用交易代替

}

//创建区块，对Block中的每一个字段填充数据即可
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{}, //先填充为空，后续填充数据
		Data:          []byte(data),
	}
	return &block

}

func main() {
	data := "hello world"
	var prevBlockHash []byte = []byte{0x0000000000000000}
	block := NewBlock(data, prevBlockHash)
	fmt.Printf("PrevBlockHash:%x\n", block.PrevBlockHash)
	fmt.Printf("Hash:%x\n", block.Hash)
	fmt.Printf("Data:%x\n", block.Data)
	fmt.Println("hello world")

}
