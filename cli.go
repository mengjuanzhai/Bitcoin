package main

import (
	"fmt"
	"os"
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
	./blockchain addBlock "xxxxxx" 添加数据到区块链
	./blockchain printBlock 打印区块链
	./blockchain getBalance 地址 获取地址的余额
`

type CLI struct {
	bc *Blockchain
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
	case "addBlock":
		if len(cmds) != 3 {
			fmt.Println(Usage)
			os.Exit(1)
		}
		fmt.Printf("添加区块链命令被调用，数据：%s\n", cmds[2])
		//data := cmds[2]
		//cli.AddBlock(data)//TODO
	case "printBlock":
		fmt.Printf("打印区块链命令被调用\n")
		cli.printBlock()
	case "getBalance":
		fmt.Printf("获取余额命令被调用\n")
		cli.bc.GetBanlance(cmds[2])

	default:
		fmt.Printf("无效的命令，请检查!\n")
		fmt.Printf(Usage)

	}

}
