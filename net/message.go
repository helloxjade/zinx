package net

import (
	iface2 "github.com/helloxjade/zinx/iface"
)

type Message struct {
	//	数据
	data []byte
	//数据长度
	len uint32
	//描述消息类型的字段
	msgid uint32
}

//创建message的方法
func NewMessage(data []byte, len, msgid uint32) iface2.IMessage {
	return &Message{
		data:  data,
		len:   len,
		msgid: msgid,
	}
}

//方法
func (msg *Message) GetData() []byte {
	return msg.data
}
func (msg *Message) GetDataLen() uint32 {
	return msg.len
}
func (msg *Message) GetMsgId() uint32 {
	return msg.msgid
}
func (msg *Message) SetData(data []byte) {
	msg.data = data
}
func (msg *Message) SetDataLen(len uint32) {
	msg.len = len
}
func (msg *Message) SetMsgId(msgid uint32) {
	msg.msgid = msgid

}
