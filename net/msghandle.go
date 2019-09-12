package net

import (
	"fmt"
	iface2 "github.com/helloxjade/zinx/iface"
	"github.com/v11-connproperity/zinx/config"
)

type MsgHandle struct {
	//维护一个map集合，key：消息id，value：路由
	//map[key] ==> value
	msghandle map[uint32]iface2.IRouter
	worksize  int                    //worker数量，由配置文件设置
	taskQueue []chan iface2.IRequest //每一个worker对应一个消息队列，每一个队列都是request的chan
}

func NewMsgHandle() *MsgHandle {

	workSize := config.GlobalConfig.WorkSize
	//别忘记make
	return &MsgHandle{
		msghandle: make(map[uint32]iface2.IRouter),
		worksize:  workSize,
		taskQueue: make([]chan iface2.IRequest, workSize),
	}
}

//1.启动worker,在服务器启动时调用，给每一个消息队列分配空间，并监听任务
func (mh *MsgHandle) StartWOrkPool() {
	fmt.Println("[StartWork Pool]...")
	//	给每一个消息队列分配空间
	for i := 0; i < mh.worksize; i++ {
		fmt.Println("启动worker,worker id:", i)
		mh.taskQueue[i] = make(chan iface2.IRequest, config.GlobalConfig.TaskQueSize)
		// 并发监听任务，每一个worker监听自己的队列
		go func(i int) {
			for {
				req := <-mh.taskQueue[i]
				fmt.Println("发现任务，执行workerid:", i)
				mh.DoMsgHandler(req)
			}
		}(i)
	}
}

//提供一个方法，向任务队列发送请求
func (mh *MsgHandle) SendReqToQueue(req iface2.IRequest) {
	//每一个链接分配一个worker，
	//同一个worker可以服务多个链接
	//1. 先获取连接cid
	cid := req.GetConnection().GetConnID()
	//	//获取当前连接所分配的workerid
	workerid := int(cid) % mh.worksize
	fmt.Println("添加cid:", cid, " 的请求到workerid:", workerid)
	//将请求放入对应的worker的消息队列中
	mh.taskQueue[workerid] <- req
}

//实现2个方法
//1.提供注册路由的方法
func (mh *MsgHandle) AddRouter(msgid uint32, router iface2.IRouter) {
	//	如果这对数据存在，则不需要添加
	_, ok := mh.msghandle[msgid]
	if ok {
		fmt.Println("路由存在，不需要添加，msgid:", msgid)
		return
	}
	//	添加路由
	mh.msghandle[msgid] = router
	fmt.Println("添加路由成功，msgid:", msgid)
}

//执行路由的handle函数
func (mh *MsgHandle) DoMsgHandler(req iface2.IRequest) {
	//	里面调用三个处理函数
	msgid := req.GetMessage().GetMsgId()
	router, ok := mh.msghandle[msgid]
	if !ok {
		fmt.Println("不存在msgid:%d 对应的路由！", msgid)
		return
	}
	router.PreHandle(req)
	router.Handle(req)
	router.PostHandle(req)
}
