package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

const GENENISISINFO = "The Times 08/Feb/2021 Chancellor on brink of second bailout for banks"
const BLOCKCHAINNAME = "blockchain.db"
const BLOCKBUCKET = "blockBucket"
const LASTHASHKEY = "lastHashkey"

//使用bolt改写
type Blockchain struct {
	//操作数据库的句柄
	db *bolt.DB
	//尾巴，存储最后一个区块的哈希
	tail []byte
}

var tail []byte

//实现创建区块链的方法
func NewBlockchain(miner string) *Blockchain {
	//功能分析：
	//1、获得数据库句柄，打开数据库，读取数据
	db, err := bolt.Open(BLOCKCHAINNAME, 0600, nil)
	if err != nil {
		log.Panic()
	}
	//defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		//判断是否有bucket，如果没有，创建bucket
		b := tx.Bucket([]byte(BLOCKBUCKET))
		if b == nil {
			fmt.Println("bucket不存在，准备创建！")
			b, err = tx.CreateBucket([]byte(BLOCKBUCKET))
			if err != nil {
				log.Panic()
			}
			//写入创世块
			//创世块中只有一个挖矿交易，只有Coinbase
			coinbase := NewCoinbaseTX(miner)

			genesisBlock := NewBlock([]*Transaction{coinbase}, []byte{})
			b.Put(genesisBlock.Hash, genesisBlock.Serialize()) //将区块链序列化，转换为字节流
			//写入lastHashKey这条数据
			b.Put([]byte(LASTHASHKEY), genesisBlock.Hash)
			tail = genesisBlock.Hash
			/*//测试
			blockInfo := b.Get(genesisBlock.Hash)
			fmt.Println("解码后的数据：", Deserialize(blockInfo))*/
		} else {
			tail = b.Get([]byte(LASTHASHKEY))
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

//添加区块
func (bc *Blockchain) addBlock(txs []*Transaction) {
	bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKBUCKET))
		if b == nil {
			os.Exit(1)
		}
		block := NewBlock(txs, bc.tail)
		b.Put(block.Hash, block.Serialize())
		b.Put([]byte(LASTHASHKEY), block.Hash)
		bc.tail = block.Hash
		return nil
	})
}

//定义区块链的迭代器，包含两个元素db,current
//db：为了遍历账本
//current：为了访问每一个区块
//迭代器要实现一个方法：Next()
//每次调用Next方法，做两件事情
//1、返回当前所指向的区块链数据block
//2、指针向前移动
type BlockchainIterator struct {
	db      *bolt.DB
	current []byte
}

//创建迭代器，要对bc进行初始化
func (bc *Blockchain) NewIterator() *BlockchainIterator {
	return &BlockchainIterator{bc.db, bc.tail}

}
func (it *BlockchainIterator) Next() *Block {
	var block Block
	err := it.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKBUCKET))
		if b == nil {
			os.Exit(1)
		}
		//真正读取数据
		blockInfo /*block的字节流*/ := b.Get(it.current)
		block = *Deserialize(blockInfo)
		it.current = block.PrevBlockHash
		return nil
	})
	if err != nil {
		log.Panic()

	}
	return &block

}

//计算余额
//实现思路

func (bc *Blockchain) FindMyUtoxs(address string) []TXoutput {
	fmt.Printf("FindMyUtoxs\n")
	//TODO
	//1、遍历账本
	it := bc.NewIterator()
	for {
		block := it.Next()
		//2、遍历交易
		for _, tx := range block.Transactions {

		}

	}

	//3、遍历output
	//4、找到属于我的所有output
	return []TXoutput{}
}
func (bc *Blockchain) GetBanlance(address string) {
	utxos := bc.FindMyUtoxs(address)
	var total = 0.0
	for _, utxo := range utxos {
		total += utxo.Value
	}
	fmt.Printf("%s的余额为：%f\n", address, total)
}

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
