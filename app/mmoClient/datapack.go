package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Message struct {
	len   uint32
	msgid uint32
	data  []byte
}
func NewMessage(msgid uint32,data []byte)*Message{
	return &Message{
		len:   uint32(len(data)),
		msgid: msgid,
		data:  data,
	}

}
type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

//封包（LTV）：L：len T:type ,V:value
//len //msgid//data
func (dp *DataPack) Pack(message *Message) ([]byte, error) {
	var buff bytes.Buffer
	//	写入数据长度
	if err := binary.Write(&buff, binary.LittleEndian, &message.len); err != nil {
		fmt.Println("binary write len err:", err)
		return nil, err
	}
	//写入数据类型
	if err := binary.Write(&buff, binary.LittleEndian, &message.msgid); err != nil {
		fmt.Println("binary write err:", err)

		return nil, err
	}

	//写入数据
	if err := binary.Write(&buff, binary.LittleEndian, &message.data); err != nil {
		fmt.Println("binary write err:", err)
		return nil, err
	}
	return buff.Bytes(), nil
}

//拆包
func (dp *DataPack) UnPack(data []byte) (*Message, error) {
	fmt.Println("开始拆包...")
	var msg Message
	//	在外面创建一个reader
	reader := bytes.NewReader(data)
	//1.     读取len
    if err:=binary.Read(reader,binary.LittleEndian,&msg.len);err!=nil{
    	fmt.Println("binary read len err:",err)
    	return nil,err
	}
	//2. 读取 msgid
	if err:=binary.Read(reader,binary.LittleEndian,&msg.msgid);err!=nil{
		fmt.Println("binary read msgid err:",err)
		return nil,err
	}
	return &msg,nil
}
