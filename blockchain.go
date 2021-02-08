package main

const GENENISISINFO = "The Times 08/Feb/2021 Chancellor on brink of second bailout for banks"

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
