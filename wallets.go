package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
)

//Wallets结构
//把地址和密钥对对应起来
//map[address1]->walletKeyPair1
//map[address2]->walletKeyPair2
//map[address3]->walletKeyPair3
//这个Wallets是对外的，WalletKeyPair是对内的
//Wallets调用WalletKeyPair
type Wallets struct {
	WalletsMap map[string]*WalletKeyPair
}

//创建wallets，返回wallets实例
func NewWallets() *Wallets {
	var ws Wallets
	ws.WalletsMap = make(map[string]*WalletKeyPair)
	//把所有的钱包从本地加载出来
	if !ws.LoadFromFile() {
		fmt.Println("加载钱包数据失败!")
	}
	//2、把实例返回
	return &ws
}

const WalletName = "wallet.dat"

func (ws *Wallets) CreateWallet() string {
	//调用NewWalletKeyPair
	wallet := NewWalletKeyPair()
	//将返回的walletKeyPair添加到WalletMap中
	address := wallet.GetAdress()
	ws.WalletsMap[address] = wallet
	//保存到本地文件
	res := ws.SaveToFile()
	if !res {
		fmt.Println("创建钱包失败！")
		return ""
	}
	return address
}
func (ws *Wallets) SaveToFile() bool {
	var buffer bytes.Buffer
	//将接口文件明确注册一下，否则gob注册失败！
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(ws)
	if err != nil {
		fmt.Printf("钱包序列化失败！，err: %v \n", err)
		return false
	}
	content := buffer.Bytes()
	//func WriteFile(filename string, data []byte, perm os.FileMode) error
	err1 := ioutil.WriteFile(WalletName, content, 0600)
	if err1 != nil {
		fmt.Printf("钱包创建失败！\n")
		return false
	}
	return true

}
func (ws *Wallets) LoadFromFile() bool {
	//读取文件
	//文件解码
	//赋值给ws
	if !IsFileExist(WalletName) {
		fmt.Println("钱包文件不存在，准备创建！")
		return true
	}
	content, err := ioutil.ReadFile(WalletName)
	if err != nil {
		return false
	}
	var wsLocal Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(content))
	err = decoder.Decode(&wsLocal)
	if err != nil {
		return false
	}
	ws.WalletsMap = wsLocal.WalletsMap
	return true

}

func (ws *Wallets) ListAddress() []string {
	//遍历ws.WalletsMap结构返回key即可
	var addresses []string
	for address, _ := range ws.WalletsMap {
		addresses = append(addresses, address)
	}
	return addresses
}
