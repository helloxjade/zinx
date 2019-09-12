package net

import (
	iface2 "github.com/helloxjade/zinx/iface"
)

type Request struct {
	conn iface2.IConnection
	//data []byte
	//len uint32
	message iface2.IMessage
}

func NewRequest(conn iface2.IConnection, msg iface2.IMessage) iface2.IRequest {
	return &Request{
		conn: conn,
		//data:data,
		//len:  len,
		message: msg,
	}
}
func (req *Request) GetConnection() iface2.IConnection {
	return req.conn
}
func (req *Request) GetMessage() iface2.IMessage {
	return req.message
}
