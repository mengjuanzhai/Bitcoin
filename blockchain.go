package main

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"github.com/boltdb/bolt"
	"lib/base58"
	"log"
	"os"
)

const GENENISISINFO = "The Times 08/Feb/2021 Chancellor on brink of second bailout for banks"
const BLOCKCHAINNAME = "blockchain.db"
const BLOCKBUCKETNAME = "blockBucket"
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
func CreateBlockchain(miner string) *Blockchain {
	if IsFileExist(BLOCKCHAINNAME) {
		fmt.Println("区块链已经存在了，不需要重复创建！")
		return nil
	}
	//功能分析：
	//1、获得数据库句柄，打开数据库，读取数据
	db, err := bolt.Open(BLOCKCHAINNAME, 0600, nil)
	if err != nil {
		log.Panic()
	}

	db.Update(func(tx *bolt.Tx) error {
		//创建bucket
		b, err := tx.CreateBucket([]byte(BLOCKBUCKETNAME))
		if err != nil {
			log.Panic()
		}
		//写入创世块
		//创世块中只有一个挖矿交易，只有Coinbase
		coinbase := NewCoinbaseTX(miner, GENENISISINFO)

		genesisBlock := NewBlock([]*Transaction{coinbase}, []byte{})
		b.Put(genesisBlock.Hash, genesisBlock.Serialize()) //将区块链序列化，转换为字节流
		//写入lastHashKey这条数据
		b.Put([]byte(LASTHASHKEY), genesisBlock.Hash)
		tail = genesisBlock.Hash
		/*//测试
		blockInfo := b.Get(genesisBlock.Hash)
		fmt.Println("解码后的数据：", Deserialize(blockInfo))*/

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

//返回区块链实例
func NewBlockchain() *Blockchain {
	if !IsFileExist(BLOCKCHAINNAME) {
		fmt.Println("区块链不存在，请先创建！")
		return nil
	}
	//功能分析：
	//1、获得数据库句柄，打开数据库，读取数据
	db, err := bolt.Open(BLOCKCHAINNAME, 0600, nil)
	if err != nil {
		log.Panic()
	}
	//defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKBUCKETNAME))
		if b == nil {
			fmt.Println("区块链bucket为空，请检查！")
			os.Exit(1)
		}
		tail = b.Get([]byte(LASTHASHKEY))
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
func (bc *Blockchain) AddBlock(txs []*Transaction) {
	//矿工得到交易时，第一时间对交易进行验证
	//矿工如果不验证，即使挖矿成功，广播区块之后，其他验证矿工，仍然会验证每一笔交易
	validTXs := []*Transaction{}
	for _, tx := range txs {
		if bc.VerifyTransaction(tx) {
			fmt.Printf("有效交易：%x\n", tx.TXid)
			validTXs = append(validTXs, tx)
		} else {
			fmt.Printf("发现无效的交易：%x\n", tx.TXid)
		}

	}
	bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKBUCKETNAME))
		if b == nil {
			os.Exit(1)
		}
		block := NewBlock(validTXs, bc.tail)
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
		b := tx.Bucket([]byte(BLOCKBUCKETNAME))
		if b == nil {
			fmt.Println("bukect为空，请检查！")
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
func (bc *Blockchain) GetBalance(address string) {
	decodeInco := base58.Decode(address)
	pubKeyHash := decodeInco[1 : len(decodeInco)-4]
	utxoinfos := bc.FindMyUtoxs(pubKeyHash)
	var total = 0.0
	for _, utxo := range utxoinfos {
		total += utxo.Output.Value
	}
	fmt.Printf("%s的余额为：%f\n", address, total)
}

type UTXOInfo struct {
	TXID   []byte   //交易ID
	Index  int64    //output的索引
	Output TXoutput //output本身
}

//计算余额
//实现思路
//1、遍历账本
//2、遍历交易
//3、遍历output
//4、找到属于我的所有output
func (bc *Blockchain) FindMyUtoxs(pubKeyHash []byte) []UTXOInfo {
	fmt.Printf("FindMyUtoxs\n")
	var UTXOInfos []UTXOInfo
	//1、遍历账本
	it := bc.NewIterator()
	//这是标识已经消耗过的utxo的机构，key是交易id，value是这个id里面的output索引的数组
	spentUTXOs := make(map[string][]int64)
	for {
		block := it.Next()
		//2、遍历交易
		for _, tx := range block.Transactions {
			//遍历输入：inputs
			//如果不是coinbase，说明是普通交易，才有必要进行遍历
			if tx.IsCoinbase() == false {
				for _, input := range tx.TXInputs {
					if bytes.Equal(HashPubKey(input.PublicKey), pubKeyHash) {
						fmt.Printf("找到已经消耗过的output！index：%d", input.Index)
						key := string(input.TXID)
						spentUTXOs[key] = append(spentUTXOs[key], input.Index)
					}
				}
			}

			//3、遍历output
		OUTPUT:
			for i, output := range tx.TXoutputs {
				key := string(tx.TXid)
				indexs := spentUTXOs[key]
				if len(indexs) != 0 {
					fmt.Printf("当前这笔交易中有被消耗过的output！\n")
					for _, j := range indexs {
						if int64(i) == j {
							fmt.Printf("i == j,当前的output已经消耗过了，跳过不统计！\n")
							continue OUTPUT

						}
					}

				}
				//4、找到属于我的所有output
				if bytes.Equal(pubKeyHash, output.PublicKeyHash) {
					//fmt.Printf("找到了属于%s的output,i : %d\n", address, i)
					utxoinfo := UTXOInfo{tx.TXid, int64(i), output}
					UTXOInfos = append(UTXOInfos, utxoinfo)
				}
			}
		}
		if len(block.PrevBlockHash) == 0 {
			fmt.Printf("遍历区块链结束！\n")
			break
		}

	}
	return UTXOInfos

}

//找到合适的utxos
func (bc *Blockchain) FindSuitableUTXOs(pubKeyHash []byte, amount float64) (map[string][]int64, float64) {
	validUTXOs := make(map[string][]int64) //标识有用的utxo
	//+++++++++++++++++++++++++++
	var resValue float64 //这些utxo存储的金额
	/*//1、遍历账本
	it := bc.NewIterator()
	//这是标识已经消耗过的utxo的机构，key是交易id，value是这个id里面的output索引的数组
	spentUTXOs := make(map[string][]int64)*/
	//复用FindMyUtoxs函数，这个函数已经包含了所有的信息
	utxoinfos := bc.FindMyUtoxs(pubKeyHash)
	for _, utxoinfo := range utxoinfos {
		key := string(utxoinfo.TXID)
		//在这里，实现了控制逻辑
		//找到符合条件的output
		//添加到返回结构validUTXOs中
		validUTXOs[key] = append(validUTXOs[key], utxoinfo.Index)
		//判断一下金额是否足够
		//a.足够，直接返回
		//b.不足，继续遍历
		resValue += utxoinfo.Output.Value
		if resValue >= amount {
			break
		}

	}

	/*for {
		block := it.Next()
		//2、遍历交易
		for _, tx := range block.Transactions {
			//遍历输入：inputs
			for _, input := range tx.TXInputs {
				if input.Address == from {
					fmt.Printf("找到已经消耗过的output！index：%d", input.Index)
					key := string(input.TXID)
					spentUTXOs[key] = append(spentUTXOs[key], input.Index)
				}
			}
		OUTPUT:
			//3、遍历output
			for i, output := range tx.TXoutputs {
				key := string(tx.TXid)
				indexs := spentUTXOs[key]
				if len(indexs) != 0 {
					fmt.Printf("当前这笔交易中有被消耗过的output！\n")
					for _, j := range indexs {
						if int64(i) == j {
							fmt.Printf("i == j,当前的output已经消耗过了，跳过不统计！\n")
							continue OUTPUT
						}
					}

				}
				//4、找到属于我的所有output
				if from == output.Address {
					fmt.Printf("找到了属于%s的output,i : %d\n", from, i)
					//UTXOs = append(UTXOs, output)
					//在这里，实现了控制逻辑
					//找到符合条件的output
					//添加到返回结构validUTXOs中
					validUTXOs[key] = append(validUTXOs[key], int64(i))
					//判断一下金额是否足够
					//a.足够，直接返回
					//b.不足，继续遍历
					resValue += output.Value
					if resValue >= amount {
						break
					}
				}
			}
		}
		if len(block.PrevBlockHash) == 0 {
			fmt.Printf("遍历区块链结束！\n")
			break
		}

	}*/

	return validUTXOs, resValue
}
func (bc *Blockchain) SignTransaction(tx *Transaction, privateKey *ecdsa.PrivateKey) {
	//1、遍历账本找到所有的交易
	prevTxs := make(map[string]Transaction)
	//遍历tx的inputs,通过id去查找所应用的交易
	for _, input := range tx.TXInputs {
		prevTx := bc.FindTransaction(input.TXID)
		if prevTx == nil {
			fmt.Printf("没有找到交易：%x\n", input.TXID)
		} else {
			//把找到的引用交易保存起来
			//0x222
			//0x333
			prevTxs[string(input.TXID)] = *prevTx
		}
	}
	tx.Sign(privateKey, prevTxs)
}
func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	//校验时，如果是挖矿交易coinbase,直接返回true
	if tx.IsCoinbase() {
		return true
	}
	//1、遍历账本找到所有的交易
	prevTxs := make(map[string]Transaction)
	//遍历tx的inputs,通过id去查找所应用的交易
	for _, input := range tx.TXInputs {
		prevTx := bc.FindTransaction(input.TXID)
		if prevTx == nil {
			fmt.Printf("没有找到交易：%x\n", input.TXID)
		} else {
			//把找到的引用交易保存起来
			//0x222
			//0x333
			prevTxs[string(input.TXID)] = *prevTx
		}
	}
	return tx.Verify(prevTxs)

}
func (bc *Blockchain) FindTransaction(txid []byte) *Transaction {
	//遍历区块链的交易
	//通过对比id来识别
	it := bc.NewIterator()
	for {
		block := it.Next()
		for _, tx := range block.Transactions {
			if bytes.Equal(tx.TXid, txid) {
				fmt.Printf("找到所引用交易：%x\n", tx.TXid)
				return tx
			}
			if len(block.PrevBlockHash) == 0 {
				break
			}
		}

	}
}

//矿工校验流程
//1、找到交易input所引用的所有交易的prevTXs
//2、对交易进行校验

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
