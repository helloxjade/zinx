package net

import (
	"fmt"
	iface2 "github.com/helloxjade/zinx/iface"
	"github.com/v11-connproperity/zinx/config"
	"io"
	"net"
	"sync"
)

//定义Connection结构体
type Connection struct {
	conn     *net.TCPConn //原生socket的链接，connection,与客户端建立链接
	connID   uint32       //每一个链接一个唯一的id
	isClosed bool         //Stop函数里使用
	//callback iface.CallBackFunc//用户注册的业务处理函数
	//router iface.IRouter
	msgHandle *MsgHandle //这里面包含了路由和msgid的集合
	msgChan   chan []byte
	exitChan  chan bool

	server    iface2.IServer //conn保存自己所属的server
	properity map[string]interface{}
	plock     sync.RWMutex
}

//给结构体赋值的过程
func NewConnection(conn *net.TCPConn, cid uint32, mh *MsgHandle, server iface2.IServer) iface2.IConnection {
	return &Connection{
		conn:     conn,
		connID:   cid,
		isClosed: false,
		//callback: callback,
		//router: router,
		msgHandle: mh,
		msgChan:   make(chan []byte),
		exitChan:  make(chan bool),
		server:    server,
		properity: make(map[string]interface{}),
	}
}
func (c *Connection) Start() {
	fmt.Println("[Connection Start]...")
	go c.StartReader()
	go c.StartWriter()
	//	调用已经注册好的钩子函数
	c.server.CallStartHook(c)
}
func (c *Connection) StartWriter() {
	fmt.Println("[StarWrite Start...]")
	defer fmt.Println("StartWriter goroution exit")
	for {
		select {
		case sendinfo := <-c.msgChan:
			_, err := c.conn.Write(sendinfo)
			if err != nil {
				fmt.Println("tcpconn.Write err:", err)
				return
			}
		//	在stop中写进true
		case <-c.exitChan:
			return
		}
	}

}

//实现接口方法
func (c *Connection) StartReader() {
	fmt.Println("[Start Reader]...")
	defer fmt.Println("StartReader goroutine exit")
	defer c.Stop()
	//处理具体的业务
	for {
		//buf:=make([]byte,512)
		//cnt,err:=c.conn.Read(buf)
		//if err!=nil{
		//	fmt.Println("tcpconn.Read.err:",err)
		//	break
		//}
		////读取的是客户端传来的数据
		//fmt.Println("Server<====Client,len:",cnt,",buf:",string(buf[:cnt]))

		//	拆包
		dp := NewDataPack()
		//1.读取8字节，解析出长度和消息id
		//先定义一个装协议头的buffer
		headBuffer := make([]byte, dp.GetDataHeadLen())

		//读取指定长度的数据
		cnt, err := io.ReadFull(c.conn, headBuffer)
		if err != nil {
			fmt.Println("io.Readfull err:", err)
			return
		}
		fmt.Printf("读取数据头的长度：%d\n", cnt)
		err, message := dp.Unpack(headBuffer)
		if err != nil {
			fmt.Printf("拆包之后的message数据详情：%v\n", message)
		}
		//检验数据包是否有效
		dataLen := message.GetDataLen()
		if dataLen == 0 {
			fmt.Printf("数据长度为0，无需读取，msgid:%d\n", message.GetMsgId())
			continue
		}
		//创建一个buf,用于存储真是的数据
		databuf := make([]byte, dataLen)
		//2.第二次读取，读取真实的数据（长度为N)
		cnt, err = io.ReadFull(c.conn, databuf)
		fmt.Printf("Server<========Client,data:%s,cnt:%d，msgid:%d\n", databuf, cnt, message.GetMsgId())
		//3.拼出满足条件的msg,赋值给request
		//把data 写进message
		message.SetData(databuf)
		//	创建request,所有的操作交给request
		req := NewRequest(c, message)
		//	具体的业务由用户传入的处理函数来执行
		//c.callback(req)
		//c.router.PreHandle(req)
		//c.router.Handle(req)
		//c.router.PostHandle(req)
		if config.GlobalConfig.WorkSize > 0 {
			go c.msgHandle.SendReqToQueue(req)
		} else {
			go c.msgHandle.DoMsgHandler(req)
		}
	}
}
func (c *Connection) Stop() {
	fmt.Println("[Connection Stop...]")
	if c.isClosed {
		return
	}
	c.isClosed = true
	//在connection停止时关闭连接
	c.server.GetConnMgr().Remove(int(c.connID))
	c.server.CallStopHookFunc(c)
	c.exitChan <- true
	//关闭chanel
	close(c.exitChan)
	close(c.msgChan)
	_ = c.conn.Close()

}

//conn提供向客户端发送数据的方法
//服务器返回的数据需要一个对应的消息id
func (c *Connection) Send(data []byte, msgid uint32) (int, error) {
	//封包
	dp := NewDataPack()
	//将服务端相应的数据返回给客户端
	err, sendinfo := dp.Pack(NewMessage(data, uint32(len(data)), msgid))
	if err != nil {
		fmt.Println("db.Pack err", err)
		return -1, err
	}
	c.msgChan <- sendinfo

	return len(sendinfo), nil
}

//获取具体的链接请求
func (c *Connection) GetConnID() uint32 {
	return c.connID
}
func (c *Connection) GetTCPConn() *net.TCPConn {
	return c.conn
}
func (c *Connection) SetProperity(key string, value interface{}) {
	c.plock.Lock()
	c.properity[key] = value
	fmt.Println("key:", key, ", value:", value)
	c.plock.Unlock()
}
func (c *Connection) GetProperity(key string) interface{} {
	fmt.Println("=======================")
	c.plock.RLock()
	value := c.properity[key]
	c.plock.RUnlock()
	fmt.Println("value:", value)
	return value
}
func (c *Connection) RemoveProperity(key string) {
	c.plock.Lock()
	delete(c.properity, key)
	c.plock.Unlock()
}
