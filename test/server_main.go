package main

import (
	"03-Zinx/v11-connproperity/zinx/iface"
	"03-Zinx/v11-connproperity/zinx/net"
	"fmt"
	"strings"
)

//1. 实现自己的路由结构：类似MainController 2. 我们的路由继承Zinx的路由：类似于MainController继承beego.Controll
//具体业务应该有=由使用框架的人传入
type TestRouter struct {
	net.Router
}

//用户重写三个函数，从而实现自己的业务
func (r *TestRouter) PreHandle(req iface.IRequest) {
	fmt.Println("用户自己实现的PreHandle")
}
func (r *TestRouter) Handle(req iface.IRequest) {
	fmt.Println("用户自己实现的Handle")
	data := req.GetMessage().GetData()
	conn := req.GetConnection() //conn提供向客户端发送数据的方法

	//	客户端发送给服务端的数据	//    转成大写返回
	writeBackInfo := strings.ToUpper(string(data))
	writeBackInfo1 := append([]byte(writeBackInfo), []byte("你好")...)
	//将回写的操作写到一个方法，由conn提供 ===> Send
	cnt, err := conn.Send(writeBackInfo1, 200)
	if err != nil {
		fmt.Println("tcpconn.Write err:", err)
		return
	}
	fmt.Println("Server=====>Client,len:", cnt, ",buf:", string(writeBackInfo1))
}
func (r *TestRouter) PostHandle(req iface.IRequest) {
	fmt.Println("用户自己实现的PostHandle")
}

//============================================
type MoveRouter struct {
	net.Router
}

func (router *MoveRouter) Handle(req iface.IRequest) {
	fmt.Println("处理移动请求的路由逻辑")
}

type Attackrouter struct {
	net.Router
}

func (router *Attackrouter) Handle(req iface.IRequest) {
	fmt.Println("处理攻击请求的路由逻辑")
}

//实现两个钩子函数
func OnConnBegin(conn iface.IConnection) {
	_, _ = conn.Send([]byte("上线成功"), 300)
	conn.SetProperity("name", "tom")
	conn.SetProperity("age", 14)
	conn.SetProperity("sex", "男")
	fmt.Println("玩家上线成功9999999999999999999")
}
func OnConnEnd(conn iface.IConnection) {
	fmt.Println("玩家下线")
	v1 := conn.GetProperity("name")
	v2 := conn.GetProperity("age")
	v3 := conn.GetProperity("sex")
	fmt.Printf("获取到的玩家的属性为：v1=%v;v2=%v,v3=%v\n", v1, v2, v3)

}
func main() {
	server := net.NewServer("zinx v.10")
	//server.AddRouter(&TestRouter{})
	server.AddRouter(1, &MoveRouter{})
	server.AddRouter(2, &Attackrouter{})

	//注册两个钩子函数
	server.RegistStartHookFunc(OnConnBegin)
	server.RegistStopHookFunc(OnConnEnd)
	server.Serve()
}
