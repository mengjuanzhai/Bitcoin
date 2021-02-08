package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

//这是一个工具类文件
func uintToByte(num uint64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	//err := binary.Write(&buffer, binary.BigEndian, num)
	//这个函数的目的是将任意的数据转换为byte字节流，这个过程叫做序列化
	//同样，可以通过binary.Read方式进行反序列化，从字节流转换为原始结构
	//binary.Read(buf,binary.LittleEndian,&num)
	//特点是：高效
	//如果在编码中有不确定长度的类型的时候，那么会报错。这时就可以使用gob来编码
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}
