package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Version       uint64
	PrevBlockHash []byte //前一哈希值
	MerkleRoot    []byte //先为空
	Timestamp     uint64
	Difficulity   uint64
	Nonce         uint64

	Hash []byte //当前哈希值,区块中不存在，为了方便我们添加进来
	Data []byte //数据，目前使用字节流，v4开始使用交易代替

}

//创建区块，对Block中的每一个字段填充数据即可
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		Version:       00,
		PrevBlockHash: prevBlockHash,
		MerkleRoot:    []byte{},
		Timestamp:     uint64(time.Now().Unix()),
		Difficulity:   BITS,
		Nonce:         10,
		Hash:          []byte{}, //先填充为空，后续填充数据
		Data:          []byte(data),
	}
	pow := NewProofOfWork(&block)
	hash, nonce := pow.Run()
	block.Nonce = nonce
	block.Hash = hash

	return &block
}

//序列化，将区块转换为字节流
func (block *Block) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()

}

//反序列化，将字节流转化为Block区块类型

func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block

}

//为了生成哈希，我们实现一个简单的函数，来计算哈希值，没有随机值，没有难度值
/*
func (block *Block) setHash() {
	var data []byte

	//data = append(data, (uintToByte(block.Version))...)
	//data = append(data, block.PrevBlockHash...)
	//data = append(data, block.MerkleRoot...)
	//data = append(data, (uintToByte(block.Timestamp))...)
	//data = append(data, (uintToByte(block.Difficulity))...)
	//data = append(data, (uintToByte(block.Nonce))...)
	//data = append(data, block.Data...)

	//使用bytes.Join方法对以上冗余代码进行优化
	dataTmp := [][]byte{
		uintToByte(block.Version),
		block.PrevBlockHash,
		block.MerkleRoot,
		uintToByte(block.Timestamp),
		uintToByte(block.Difficulity),
		uintToByte(block.Nonce),
		block.Data,
	}
	data = bytes.Join(dataTmp, []byte{})
	hash := sha256.Sum256(data)
	block.Hash = hash[:]
}
*/
