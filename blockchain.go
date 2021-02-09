package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

const GENENISISINFO = "The Times 08/Feb/2021 Chancellor on brink of second bailout for banks"

/*
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

//添加区块
func (bc *Blockchain) addBlock(data string) {
	lastBlock := bc.blocks[len(bc.blocks)-1]
	prevBlockHash := lastBlock.Hash
	block := NewBlock(data, prevBlockHash)
	bc.blocks = append(bc.blocks, block)
}
*/
//使用bolt改写
type Blockchain struct {
	//操作数据库的句柄
	db *bolt.DB
	//尾巴，存储最后一个区块的哈希
	tail []byte
}

var tail []byte

//实现创建区块链的方法
func NewBlockchain() *Blockchain {
	//功能分析：
	//1、获得数据库句柄，打开数据库，读取数据
	db, err := bolt.Open("blockchain.db", 0600, nil)
	if err != nil {
		log.Panic()
	}
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		//判断是否有bucket，如果没有，创建bucket
		b := tx.Bucket([]byte("blockBucket"))
		if b == nil {
			fmt.Println("bucket不存在，准备创建！")
			b, err = tx.CreateBucket([]byte("blockBucket"))
			if err != nil {
				log.Panic()
			}
			//写入创世块
			genesisBlock := NewBlock(GENENISISINFO, []byte{})
			b.Put(genesisBlock.Hash, genesisBlock.toBytes()) //将区块链序列化，转换为字节流
			//写入lastHashKey这条数据
			b.Put([]byte("lastHashKey"), genesisBlock.Hash)
			tail = genesisBlock.Hash
		} else {
			tail = b.Get([]byte("lastHashKey"))
		}

		return nil
	})

	//写入创世块
	//写入lastHashKey这条数据
	//更新tail为最后一个区块的哈希
	//返回bc实例
	//2、获取最后一个区块的哈希值
	//填充给tail
	//返回实例
	return &Blockchain{db, tail}

}
