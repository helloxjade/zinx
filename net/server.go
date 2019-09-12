package net

import (
	"fmt"
	iface2 "github.com/helloxjade/zinx/iface"
	"github.com/v11-connproperity/zinx/config"
	"net"
)

//server 结构需要内容
type Server struct {
	//属性：
	//1.IP
	IP string
	//2.port
	Port uint32
	//3.name服务的名字
	Name string
	//4.port
	Version string //tcp4 tcp6
	//Router  iface.IRouter//单路由
	msghandle *MsgHandle //这里面包含了路由与msgid的集合
	//	连接管理模块
	connmgr iface2.IconnManager
	//	钩子函数，由服务器开发者提供具体业务，在客户端建立、关闭连接前，
	//  主动调用钩子函数,这两个变量用于接收注册的方法，为便后续调用
	onConnStartFunc func(conn iface2.IConnection)
	onConnStopFunc  func(conn iface2.IConnection)
}

func NewServer(name string) iface2.IServer {
	return &Server{
		IP:      config.GlobalConfig.IP, //监听所有端口
		Port:    config.GlobalConfig.Port,
		Name:    config.GlobalConfig.Name,
		Version: config.GlobalConfig.Version,
		//Router:  &Router{},
		msghandle: NewMsgHandle(),
		connmgr:   NewConnManager(), //初始化
	}
}

func (s *Server) Start() {
	fmt.Println("[Server Start]...")
	//	socket 监听
	//    l:=net.Listen("tcp",8888)
	//    c:=l.Accept()
	//    c.Read()
	address := fmt.Sprintf("%s:%d", s.IP, s.Port)
	tcpAddr, err := net.ResolveTCPAddr(s.Version, address)
	if err != nil {
		fmt.Println("Server Start err:", err)
		return
	}
	listener, err := net.ListenTCP(s.Version, tcpAddr)
	if err != nil {
		fmt.Println("net.ListenTCP err:", err)
		return
	}
	//在server启动时，将worker池启动
	s.msghandle.StartWOrkPool()

	//连接id,每创建一个新的连接，cid加1
	var cid uint32
	cid = 0
	go func() {
		for { //监听
			tcpconn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("listener.Accept err :", err)
				return
			}
			//控制最大连接数,获取当前所有的链接总数,如果已经等于最大连接数上限
			if s.connmgr.GetConnCount() == config.GlobalConfig.ConnAllowSize {
				fmt.Println("已经到达最大连接上限，当前连接被拒绝：", cid)
				//	关闭连接
				_ = tcpconn.Close()
				continue
			}
			//1. 得到tcpconn，封装自己Connection
			myconnection := NewConnection(tcpconn, cid, s.msghandle, s)
			//连接管理之添加连接
			s.connmgr.AddConn(int(cid), myconnection)
			cid++
			//2. 启动conn.start，
			//server只负责管理连接，具体的业务处理，由conn负责
			go myconnection.Start()
		}
	}()

	//buf := make([]byte, 512)
	//	cnt, err := tcpconn.Read(buf)
	//	if err != nil {
	//		fmt.Println("tcpconn.Read err:", err)
	//		return
	//	}
	//	fmt.Println("Server<====Client,len:", cnt, ",buf:", string(buf[:cnt]))

}

func (s *Server) Stop() {
	fmt.Println("[server stop]...")
}
func (s *Server) Serve() {
	fmt.Println("[Server server]...")
	s.Start()
	fmt.Println("+++++++go start done...")
	for {
	}
	//select {}
}

//func (s *Server)AddRouter(router iface.IRouter)  {
//	//接受这个路由
//	s.Router=router
//}
//Server的AddRouter函数需要修改，调用MsgHandler的AddRouter即可（注意需要增加一个msgid字 段)
func (s *Server) AddRouter(msgid uint32, router iface2.IRouter) {
	//接收这个路由
	s.msghandle.AddRouter(msgid, router)
}
func (s *Server) GetConnMgr() iface2.IconnManager {
	return s.connmgr
}

//注册两个钩子函数方法，就把钩子函数存进 serve的回调函数的字段
//之所以封装一下，是为了方便其他模块调用，并且可以提供参数校验是否为nil，也可以方便后续扩展
func (s *Server) RegistStartHookFunc(f func(connection iface2.IConnection)) {
	s.onConnStartFunc = f
}
func (s *Server) RegistStopHookFunc(f func(connection iface2.IConnection)) {
	s.onConnStopFunc = f
}
func (s *Server) CallStartHook(conn iface2.IConnection) {
	//如果用户并没有调用钩子函数
	if s.onConnStartFunc == nil {
		return
	}
	s.onConnStartFunc(conn)
}
func (s *Server) CallStopHookFunc(conn iface2.IConnection) {
	if s.onConnStopFunc == nil {
		return
	}
	s.onConnStopFunc(conn)
}
