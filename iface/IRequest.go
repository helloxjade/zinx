package iface

type IRequest interface {
	GetConnection() IConnection
	GetMessage() IMessage
}
