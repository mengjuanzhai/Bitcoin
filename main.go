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
	bc := NewBlockchain("miner")
	defer bc.db.Close()
	cil := CLI{bc}
	cil.Run()
	//bc.addBlock("second block")
	//bc.addBlock("third block")
	//bc.addBlock("four block")
	//bc.addBlock("five block")
	//it := bc.NewIterator()
	//for {
	//	block := it.Next()
	//	fmt.Println("-------------------------------")
	//	fmt.Printf("Version:%x\n", block.Version)
	//	fmt.Printf("PrevBlockHash:%x\n", block.PrevBlockHash)
	//	fmt.Printf("MerkleRoot:%x\n", block.MerkleRoot)
	//	//优化时间打印，定义时间显示格式
	//	timeFormat := time.Unix(int64(block.Timestamp), 0).Format("2006-01-02 15:01:05")
	//	fmt.Printf("Timestamp:%s\n", timeFormat)
	//	fmt.Printf("Difficulity:%x\n", block.Difficulity)
	//	fmt.Printf("Nonce:%x\n", block.Nonce)
	//	fmt.Printf("Hash:%x\n", block.Hash)
	//	fmt.Printf("Data:%s\n", block.Data)
	//	pow := NewProofOfWork(block)
	//	fmt.Printf("IsVaild:%v\n", pow.IsValid())
	//	if bytes.Equal(block.PrevBlockHash, []byte{}) {
	//		fmt.Println("区块链遍历结束！")
	//		break
	//	}
	//}

	/*bolt自带Foreach方法，迭代打印
	err := bc.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(BLOCKBUCKET))

		b.ForEach(func(k, v []byte) error {
			fmt.Printf("key=%v, value=%v\n", k, v)
			return nil
		})
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	*/
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
