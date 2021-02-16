package main

import (
	"fmt"
	"os"
	"strconv"
)

//使用命令行分析
//1、所有的支配动作交给命令行来做
//2、主函数只需调用调用命令行结构即可
//3、根据输入的不同命令，命令行做响应动作
//3.1、addBlock
//3.2、printBlock
//
//CLI:command line的缩写
//type CLI struct {
//	bc *Blockchain
//}
//
//添加区块链时：bc.addBlock(data)，data通过os.Args传入
//打印区块链：遍历区块链，不需要传入外部数据

const Usage = `
	./blockchain createBlockchain  创建区块链
	./blockchain printBlock 打印区块链
	./blockchain getBalance 地址 获取地址的余额
	./blockchain send FROM TO AMOUNT MINER  DATA 转账命令
`

type CLI struct {
	//bc *Blockchain
	//CLI不需要保存区块链实例了，所有的命令在自己调用之前，自己获取区块链实例
}

//给CLI提供一个方法，进行命令解析，从而执行调度
func (cli *CLI) Run() {
	cmds := os.Args
	if len(cmds) < 2 {
		fmt.Printf("无效的命令，请检查!\n")
		fmt.Printf(Usage)
		os.Exit(1)
	}
	switch cmds[1] {
	case "createBlockchain":
		if len(cmds) != 3 {
			fmt.Println(Usage)
			os.Exit(1)
		}
		fmt.Printf("创建区块链命令被调用!\n")
		miner := cmds[2]
		cli.CreateBlockchain(miner)
	case "printBlock":
		fmt.Printf("打印区块链命令被调用\n")
		cli.printBlock()
	case "getBalance":
		fmt.Printf("获取余额命令被调用\n")
		cli.GetBalance(cmds[2])
	case "send":
		fmt.Printf("转账命令被调用\n")
		if len(cmds) != 7 {
			fmt.Printf("send命令发现无效参数，请检查！\n")
			fmt.Println(Usage)
			os.Exit(1)

		}
		from := cmds[2]
		to := cmds[3]
		amount, _ := strconv.ParseFloat(cmds[4], 64)
		miner := cmds[5]
		data := cmds[6]
		cli.Send(from, to, amount, miner, data)
	default:
		fmt.Printf("无效的命令，请检查!\n")
		fmt.Printf(Usage)

	}

}
