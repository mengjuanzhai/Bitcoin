package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"lib/base58"
	"log"
	"math/big"
	"os"
	"strings"
)

//交易输入
//指明交易发起人可支付资金的来源，包含：
//1、引用UTXO所在交易的ID
//2、所消费UTXO在output中的索引
//3、解锁脚本
type TXInput struct {
	//引用utxo所在交易ID（知道在那个房间）
	TXID []byte
	//所消费utxo在output的索引值（具体位置）
	Index int64
	//解锁脚本（签名，公钥）
	//签名
	Signature []byte
	//公钥本身，不是公钥哈希
	PublicKey []byte
}

//交易输出
//包含资金接收方的相关信息，包含：
//接收金额（数字）
//锁定脚本（对方公钥的哈希，这个哈希可以通过地址反推出来，所以转账知道地址即可）
type TXoutput struct {
	//接受的金额
	Value float64
	//Address string
	//公钥哈希，不是公钥本身
	PublicKeyHash []byte
}

//给定转账地址，得到这的地址的公钥哈希，完成对output的锁定
func (output *TXoutput) Lock(address string) {
	//address -> public key hash
	//25字节
	decodeInfo := base58.Decode(address)
	pubKeyHash := decodeInfo[1 : len(decodeInfo)-4]
	output.PublicKeyHash = pubKeyHash
}
func NewTXoutput(value float64, address string) TXoutput {
	output := TXoutput{Value: value}
	output.Lock(address)
	return output
}

type Transaction struct {
	TXid []byte //交易ID
	//所有的inputs
	TXInputs []TXInput
	//所有的outputs
	TXoutputs []TXoutput
}

//交易ID
//一般是交易结构的哈希值（参考block序列化）
func (tx *Transaction) SetTXID() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(buffer.Bytes())
	tx.TXid = hash[:]
}

//实现挖矿交易
//特点：只有输出，没有有效的输入（不需要TXInput,包括ID，索引，签名等）
//把挖矿的人传递进来，因为有奖励
const REWARD = 12.5

func NewCoinbaseTX(address string, data string) *Transaction {
	//比特币在挖矿时，对于这个input的id填0，对索引值填0xffff，data由矿工填写，一般填所在矿池的名字
	input := TXInput{nil, -1, nil, []byte(data)}
	output := NewTXoutput(REWARD, address)
	inputs := []TXInput{input}
	outputs := []TXoutput{output}
	tx := Transaction{nil, inputs, outputs}
	tx.SetTXID()
	return &tx
}
func (tx *Transaction) IsCoinbase() bool {
	//特点：1、只有1个input 2、引用的id是nil 3、引用的索引是-1
	inputs := tx.TXInputs
	if len(inputs) == 1 && inputs[0].TXID == nil && inputs[0].Index == -1 {
		return true
	}
	return false

}

//参数：
//1、付款人
//2、收款人
//3、转账金额
//4、bc
//创建普通交易的内部逻辑
//1、遍历账本，找到属于付款人的合适的金额，把这个outputs找到
//2、如果找到的outputs的钱不足以转账，创建交易失败
//3、将outputs转成inputs
//4、创建输出，创建一个属于收款人的output
//5、如果有零钱，创建属于付款人的output
//6、设置交易ID
//7、返回交易结构
func NewTransaction(from string, to string, amount float64, bc *Blockchain) *Transaction {
	//1、打开钱包
	ws := NewWallets()
	//获取密钥对
	wallet := ws.WalletsMap[from]
	if wallet == nil {
		fmt.Printf("%s的私钥不存在，交易创建失败！\n", from)
		return nil
	}
	//2、获取公钥私钥
	privateKey := wallet.PrivateKey
	publicKey := wallet.PublicKey
	pubKeyHash := HashPubKey(publicKey)

	//map[string][]int64
	//1、遍历账本，找到属于付款人的合适的金额，把这个outputs找到
	validUTXOs := make(map[string][]int64) //标识有用的utxo
	var resValue float64                   //这些utxo存储的金额
	//第一部分，找到所需要的UTXO的集合
	validUTXOs /*本次支付所需要的utxo集合*/, resValue /*返回utxos所包含的金额*/ = bc.FindSuitableUTXOs(pubKeyHash, amount)
	//2、如果找到的outputs的钱不足以转账，创建交易失败
	if resValue < amount {
		fmt.Printf("余额不足，交易失败！！！\n")
		os.Exit(1)
	}
	////3、将outputs转成inputs
	var inputs []TXInput
	var outputs []TXoutput
	//第二部分，input的创建
	for txid, indexes /*0x333*/ := range validUTXOs {
		for _, i /*索引：0,1*/ := range indexes {
			input := TXInput{[]byte(txid), i, nil, pubKeyHash}
			inputs = append(inputs, input)
		}
	}
	//4、创建输出，创建一个属于收款人的output
	//output := TXoutput{amount, to}
	output := NewTXoutput(amount, to)
	outputs = append(outputs, output)
	//5、如果有零钱，创建属于付款人的output
	if resValue > amount {
		//output1 := TXoutput{resValue - amount, from}
		output1 := NewTXoutput(resValue-amount, from)
		outputs = append(outputs, output1)
	}

	//创建交易
	tx := Transaction{nil, inputs, outputs}
	//6、设置交易ID
	tx.SetTXID()
	//我们付款人在创建交易时，已经得到所有output的所有详细信息。
	//但是我们不去使用，因为矿工在校验时，矿工是没有这部分信息的，矿工需要遍历账本所有的引用交易
	//我们为了统一操作，所以再查询一次，进行签名。
	bc.SignTransaction(&tx, privateKey)
	//7、返回交易结构
	return &tx
}

//第一个参数是私钥
//第二个参数是这个交易input所引用的所有交易
func (tx *Transaction) Sign(privKey *ecdsa.PrivateKey, privTXs map[string]Transaction) {
	fmt.Printf("对交易进行签名...\n")
	if tx.IsCoinbase() {
		return
	}
	//1、拷贝一份交易，做相应裁剪：把每一个input的Sig和pubKey设置为nil,output不做改变
	txCopy := tx.TrimmedCopy()
	//2、遍历txCopy.inputs,把这个input所引用的output的公钥哈希拿过来，赋值给input的pubKey
	for i, input := range txCopy.TXInputs {
		//找到所引用的交易
		prevTX := privTXs[string(input.TXID)]
		output := prevTX.TXoutputs[input.Index]
		//for循环迭代出来的数据是一个副本，对这个数据input进行修改，不会对原始数据txCopy产生影响
		//input.PublicKey = output.PublicKeyHash
		//所以我们采用下标方式进行修改
		txCopy.TXInputs[i].PublicKey = output.PublicKeyHash
		//3、生成要签名的数据
		//签名就是要对数据的hash进行签名
		//我们的数据都在交易中，我们要求交易的哈希
		//Transaction的SetTXID函数就是对交易的哈希
		//我们可以使用交易的ID作为我们签名的内容
		txCopy.SetTXID()
		SignData := txCopy.TXid
		//清理
		//input.PublicKey = nil
		txCopy.TXInputs[i].PublicKey = nil
		fmt.Printf("要签名的数据， SignData:%x\n", SignData)
		//4、对数据进行签名r,s
		r, s, err := ecdsa.Sign(rand.Reader, privKey, SignData)
		//5、拼接r,s为字节流，赋值给原始的交易的Signature字段
		if err != nil {
			fmt.Printf("交易签名失败，err : %v\n", err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		tx.TXInputs[i].Signature = signature
	}

}

//拷贝一份交易，做相应裁剪：把每一个input的Sig和pubKey设置为nil,output不做改变
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXoutput
	for _, input := range tx.TXInputs {
		input1 := TXInput{input.TXID, input.Index, nil, nil}
		inputs = append(inputs, input1)
	}
	outputs = tx.TXoutputs
	tx1 := Transaction{tx.TXid, inputs, outputs}
	return tx1
}
func (tx *Transaction) Verify(privTXs map[string]Transaction) bool {
	fmt.Printf("对交易进行校验...\n")
	//1、拷贝修建的副本、
	txCopy := tx.TrimmedCopy()
	//2、遍历原始交易（不是txCopy)
	for i, input := range tx.TXInputs {
		//3、遍历原始交易的input所引用的前交易prevTX
		prevTX := privTXs[string(input.TXID)]
		output := prevTX.TXoutputs[input.Index]
		//4、找到output的公钥哈希，赋值给这个input
		txCopy.TXInputs[i].PublicKey = output.PublicKeyHash
		//5、还原签名的数据
		txCopy.SetTXID()
		verifyData := txCopy.TXid
		//清理动作，重要！
		txCopy.TXInputs[i].PublicKey = nil
		fmt.Printf("要校验的数据， VerifyData:%x\n", verifyData)
		//6、校验
		//还原签名为r,s
		signature := input.Signature
		r := big.Int{}
		s := big.Int{}
		rData := signature[:len(signature)/2]
		sData := signature[len(signature)/2:]
		r.SetBytes(rData)
		s.SetBytes(sData)
		//还原公钥为curve,X,Y
		//type PublicKey struct {
		//	elliptic.Curve
		//	X, Y *big.Int
		//}
		//公钥字节流
		pubKeyBytes := input.PublicKey
		x := big.Int{}
		y := big.Int{}
		xData := signature[:len(pubKeyBytes)/2]
		yData := signature[len(pubKeyBytes)/2:]
		x.SetBytes(xData)
		y.SetBytes(yData)
		curve := elliptic.P256()
		pulicKey := ecdsa.PublicKey{curve, &x, &y}
		//数据，签名，公钥准备完毕，开始校验
		if !ecdsa.Verify(&pulicKey, verifyData, &r, &s) {
			return false
		}
	}
	return true
}

func (tx *Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.TXid))

	for i, input := range tx.TXInputs {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.TXID))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Index))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PublicKey))
	}

	for i, output := range tx.TXoutputs {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %f", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PublicKeyHash))
	}

	return strings.Join(lines, "\n")
}
