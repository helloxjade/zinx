package net

import (
	"fmt"
	iface2 "github.com/helloxjade/zinx/iface"
	"sync"
)

type ConnManager struct {
	conns    map[int]iface2.IConnection //链接的集合
	connLock sync.RWMutex               //读写锁，用于控制操作map
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		conns:    make(map[int]iface2.IConnection),
		connLock: sync.RWMutex{},
	}
}

//增加链接
func (cm *ConnManager) AddConn(cid int, conn iface2.IConnection) {
	fmt.Println("增加新连接:", cid)
	//    加锁//多个goroutine 操作map时，一定要枷锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	//    判断当前id 的链接是否存在，如不存在，则不需要添加
	if _, ok := cm.conns[cid]; ok {
		fmt.Println("当前链接已经存在呢，无需添加", cid)
		return
	}
	cm.conns[cid] = conn
	fmt.Println("connection addto connmanage successufl id=", cid)

}

//删除链接 根据connid删除
func (cm *ConnManager) Remove(cid int) {
	fmt.Println("删除连接:", cid)
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	delete(cm.conns, cid)
}

//根据给顶cid,返回链接句柄
func (cm *ConnManager) GetConn(cid int) iface2.IConnection {
	fmt.Println("获取链接", cid)
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	return cm.conns[cid]

}

//获取当前所有的链接的总数
func (cm *ConnManager) GetConnCount() int {
	fmt.Println("获取当前链接的总数")
	return len(cm.conns)
}

//清楚所有链接
func (cm *ConnManager) ClearConn() {
	fmt.Println("清除所有的链接")
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	//   1.关闭所有的conn
	for cid, conn := range cm.conns {
		conn.Stop()
		//2.把map清空
		delete(cm.conns, cid)
	}
}
