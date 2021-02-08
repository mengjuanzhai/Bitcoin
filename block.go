package main

import "crypto/sha256"

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

//为了生成哈希，我们实现一个简单的函数，来计算哈希值，没有随机值，没有难度值
func (block *Block) setHash() {
	data := []byte{}
	data = append(data, block.PrevBlockHash...)
	data = append(data, block.Data...)
	hash := sha256.Sum256(data)
	block.Hash = hash[:]
}
