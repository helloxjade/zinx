package net

import (
	"bytes"
	"encoding/binary"
	iface2 "github.com/helloxjade/zinx/iface"
)

//负责封包与拆包
type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

//直接返回自定义的协议头的长度
func (dp *DataPack) GetDataHeadLen() uint32 {
	return 8
}

//封包函数
func (dp *DataPack) Pack(msg iface2.IMessage) (error, []byte) {
	//先获取消息内容。长度。消息id
	data := msg.GetData()
	len1 := msg.GetDataLen()
	Msgid := msg.GetMsgId()
	var buf bytes.Buffer
	//写消息长度
	err := binary.Write(&buf, binary.LittleEndian, len1)
	if err != nil {
		return err, nil
	}
	//     写消息id
	err = binary.Write(&buf, binary.LittleEndian, Msgid)
	if err != nil {
		return err, nil
	}
	//写消息体
	err = binary.Write(&buf, binary.LittleEndian, data)
	if err != nil {
		return err, nil
	}
	// 把buf类调用bytes（）方法 换成字节流
	return nil, buf.Bytes()
}

//拆包函数
func (dp *DataPack) Unpack(data []byte) (error, iface2.IMessage) {
	//	创建一个reader,读取data
	reader := bytes.NewReader(data)
	//    创建一个message结构体，用于存储拆包后的数据
	var message Message
	//拆包主要是解析出两个内容：真实传递数据的长度，消息id

	//func Read(r io.Reader, order ByteOrder, data interface{}) error
	//err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &message.len) //错误做法:
	//这样每次只能读出数据长度
	//1. 读取数据头, 获得数据长度
	err := binary.Read(reader, binary.LittleEndian, &message.len) //这样就会把数据写进message的len字段
	if err != nil {
		return err, nil
	}

	//2.读取数据头，获取数据id
	err = binary.Read(reader, binary.LittleEndian, &message.msgid)
	if err != nil {
		return err, nil
	}
	return nil, &message
}
