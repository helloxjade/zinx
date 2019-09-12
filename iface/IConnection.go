package iface

import "net"

//

//1. Start()  ===> 读，写方法  , 工作
//1. Send() ===> 向conn发送数据
//2. Stop()
//3. GetConnID() ==> 每一个连接有自己的id
//4. GetTCPConn() *netTCPConn
type IConnection interface {
	Start()                           //开启链接
	Stop()                            //关闭连接
	Send([]byte, uint32) (int, error) //把数据返回
	GetConnID() uint32                //每次请求都不一样，要知道是处理哪个请求就要获取ID
	GetTCPConn() *net.TCPConn         //原生的链接 通过这个与框架绑定
	SetProperity(string, interface{})
	GetProperity(string) interface{}
	RemoveProperity(string)
}

//定义一个回调函数，由用户提供，处理用户指定的业务
//路由
//先把这个准备好，回头再调用，用的时候再调用，注册进来时不立刻调用。
type CallBackFunc func(IRequest)
