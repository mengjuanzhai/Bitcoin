package main

//1、定义结构，区块链的字段比正常少
//a、当前哈希值
//b、前一区块的哈希值
//3、数据
//2、创建区块
//3、生成哈希
//4、引入区块链
//5、添加区块
//6、重构代码

func main() {
	bc := NewBlockchain()
	bc.db.Close()
	/*bc.addBlock("第二个区块")
	for i, block := range bc.blocks {
		fmt.Printf("--------------+%d+------------------\n", i)
		fmt.Printf("Version:%x\n", block.Version)
		fmt.Printf("PrevBlockHash:%x\n", block.PrevBlockHash)
		fmt.Printf("MerkleRoot:%x\n", block.MerkleRoot)
		//优化时间打印，定义时间显示格式
		timeFormat := time.Unix(int64(block.Timestamp), 0).Format("2006-01-02 15:01:05")
		fmt.Printf("Timestamp:%s\n", timeFormat)
		fmt.Printf("Difficulity:%x\n", block.Difficulity)
		fmt.Printf("Nonce:%x\n", block.Nonce)
		fmt.Printf("Hash:%x\n", block.Hash)
		fmt.Printf("Data:%s\n", block.Data)
		pow := NewProofOfWork(block)
		fmt.Printf("IsVaild:%v\n", pow.IsValid())
	}*/

}
