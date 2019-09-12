package iface

type IconnManager interface {
	AddConn(int, IConnection) //增加链接
	Remove(int)               //删除链接
	GetConn(int) IConnection  //获取当前链接
	GetConnCount() int        //获取当前所有的链接总数
	ClearConn()               //清楚所有的链接
}
