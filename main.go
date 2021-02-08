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
const GENENISISINFO = "The Times 08/Feb/2021 Chancellor on brink of second bailout for banks"

func main() {
	bc := NewBlockchain()
	bc.addBlock("第二个区块")
	for i, block := range bc.blocks {
		fmt.Printf("--------------+%d+------------------\n", i)
		fmt.Printf("PrevBlockHash:%x\n", block.PrevBlockHash)
		fmt.Printf("Hash:%x\n", block.Hash)
		fmt.Printf("Data:%x\n", string(block.Data))
	}

}
