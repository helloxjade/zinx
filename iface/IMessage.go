package iface

type IMessage interface {
	GetData() []byte
	GetDataLen() uint32
	GetMsgId() uint32
	SetData(data []byte)
	SetDataLen(len uint32)
	SetMsgId(msgid uint32)
}
