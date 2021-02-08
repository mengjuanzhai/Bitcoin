package main

import (
	"crypto/sha256"
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
const GENENISISINFO = "The Times 08/Feb/2021 Chancellor on brink of second bailout for banks"

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
	block.setHash()
	return &block
}

//创建区块链，使用数组进行模拟
type Blockchain struct {
	blocks []*Block
}

//实现创建区块链的方法
func NewBlockchain() *Blockchain {
	genesisBlock := NewBlock(GENENISISINFO, []byte{0x000000000000})
	bc := Blockchain{blocks: []*Block{genesisBlock}}
	return &bc

}

//为了生成哈希，我们实现一个简单的函数，来计算哈希值，没有随机值，没有难度值
func (block *Block) setHash() {
	data := []byte{}
	data = append(data, block.PrevBlockHash...)
	data = append(data, block.Data...)
	hash := sha256.Sum256(data)
	block.Hash = hash[:]
}

//添加区块
func (bc *Blockchain) addBlock(data string) {
	lastBlock := bc.blocks[len(bc.blocks)-1]
	prevBlockHash := lastBlock.Hash
	block := NewBlock(data, prevBlockHash)
	bc.blocks = append(bc.blocks, block)
}

func main() {
	bc := NewBlockchain()
	bc.addBlock("第二个区块")
	for i, block := range bc.blocks {
		fmt.Printf("--------------+%d+------------------\n", i)
		fmt.Printf("PrevBlockHash:%x\n", block.PrevBlockHash)
		fmt.Printf("Hash:%x\n", block.Hash)
		fmt.Printf("Data:%x\n", block.Data)
	}

}
