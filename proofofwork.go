package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//1、定义一个工作量证明结构ProofOfWork
//a.block
//b.目标值
//2、提供创建POW的方法（函数）：NewProofOfWork(参数)
//3、提供不断计算hash的函数：Run()
//4、提供一个校验函数：IsValid()

type ProofOfWork struct {
	block *Block
	//使用大数类big.Int来存储哈希值，它内置了一些方法
	//Cmp:比较方法
	//SetBytes : 把bytes转成big.int类型
	//SetString : 把string转成big.int类型
	target *big.Int //系统提供的，是固定的
}

const BITS = 16

func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}
	//写难度值，难度值应该是推导出来的，但是我们为了简化，把难度值写成固定的，一切完成后再推导
	//powStr := "0001000000000000000000000000000000000000000000000000000000000000"
	//固定难度值
	//16进制格式的字符串
	//var bigIntTmp big.Int
	//bigIntTmp.SetString(powStr, 16)
	//pow.target = &bigIntTmp
	bigIntTmp := big.NewInt(1)
	bigIntTmp.Lsh(bigIntTmp, 256-BITS)
	pow.target = bigIntTmp
	return &pow
}

//这是pow的运算函数，为了获取挖矿的随机数，同时返回区块的哈希值
func (pow *ProofOfWork) Run() ([]byte, uint64) {
	var nonce uint64
	var hash [32]byte
	//1、获取block数据
	//2、拼接nonce
	//3、sha256
	//4、与难度进行比较
	//a、哈希函数大于难度值，nonce++
	//b、哈希函数小于难度值，挖矿成功，返回哈希值及随机数
	for {
		data := pow.preparedData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("%x\r", hash)
		//将hash（[]byte类型）转为big.Int，然后与pow.target进行比较，需要引入局部变量
		var bigIntTemp big.Int
		bigIntTemp.SetBytes(hash[:])
		res := bigIntTemp.Cmp(pow.target)
		if res == -1 {
			//此时x<y，表示挖矿成功
			fmt.Printf("挖矿成功！nonce = %x,hash = %x\n", nonce, hash)
			break
		} else {
			nonce++
		}
	}
	return hash[:], nonce

}

func (pow *ProofOfWork) preparedData(nonce uint64) []byte {
	block := pow.block
	dataTmp := [][]byte{
		uintToByte(block.Version),
		block.PrevBlockHash,
		block.MerkleRoot,
		uintToByte(block.Timestamp),
		uintToByte(block.Difficulity),
		uintToByte(nonce),
		//block.Data
	}
	data := bytes.Join(dataTmp, []byte{})
	return data
}

//更正：比特币取哈希值，并不是对整个区块取哈希，而是对区块头取哈希

//IsValid 校验函数:校验一下，Hash，block数据和Nonce是否满足难度值要求
func (pow *ProofOfWork) IsValid() bool {
	data := pow.preparedData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	var tmpBigInt big.Int
	tmpBigInt.SetBytes(hash[:])
	return tmpBigInt.Cmp(pow.target) == -1
	//true 表示想x<y
}
