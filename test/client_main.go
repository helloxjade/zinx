package main

import (
	net2 "03-Zinx/v11-connproperity/zinx/net"
	"fmt"
	"io"
	"net"
	"time"
)

func main01() {
	conn, err := net.Dial("tcp", ":8848")
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return
	}
	data := []byte("hello world")
	for {
		cnt, err := conn.Write(data)
		if err != nil {
			fmt.Println("conn write err", err)
			continue
		}
		fmt.Println("Client===>Server,len:", cnt, "data:", data)
		buf := make([]byte, 512)
		cnt, err = conn.Read(buf)
		if err != nil {
			fmt.Println("conn read err:", err)
			continue
		}
		time.Sleep(1 * time.Second)
	}
}
func main() {
	data0 := []byte{}
	data1 := []byte("你好")
	data2 := []byte("hello world")
	data3 := []byte("国庆即将到来")
	//	创建message
	msg0 := net2.NewMessage(data0, uint32(len(data0)), 0)
	msg1 := net2.NewMessage(data1, uint32(len(data1)), 0)
	msg2 := net2.NewMessage(data2, uint32(len(data2)), 1)
	msg3 := net2.NewMessage(data3, uint32(len(data3)), 2)
	//   对message 进行封包
	dp := net2.NewDataPack()
	_, info0 := dp.Pack(msg0)
	_, info1 := dp.Pack(msg1)
	_, info2 := dp.Pack(msg2)
	_, info3 := dp.Pack(msg3)
	//将三个消息的字节流拼接到一起，一次性发送给服务器
	infosend := append(info0, info1...)
	infosend = append(infosend, info2...)
	infosend = append(infosend, info3...)
	//	发送
	conn, err := net.Dial("tcp", "127.0.0.1:8848")
	if err != nil {
		fmt.Printf("net Dial err :%v\n", err)
		return
	}
	go func() {
		for {
			fmt.Printf("++++++++++++++++++\n")
			cnt, err := conn.Write(infosend)
			if err != nil {
				fmt.Printf("Client send err:%v\n", err)
				return
			}
			fmt.Println("Client=====》Server cnt:", cnt)
			time.Sleep(3 * time.Second)
		}
	}()
	go func() {
		for {
			//客户端解析（拆包），得到服务器的相应数据
			//++++++++++++++++++++++++++++++++++++++++++++++
			//第一次解析数据头
			headbuf := make([]byte, 8)
			cnt, err := io.ReadFull(conn, headbuf)
			if err != nil {
				fmt.Println("io.ReadFull err", err)
				return
			}
			fmt.Printf("读取到的数据长度为%d\n", cnt)
			//拆包
			dp := net2.NewDataPack()
			err, message := dp.Unpack(headbuf)
			if err != nil {
				fmt.Println("db.Unpack err", err)
				return
			}
			fmt.Printf("拆包之后的数据为message:%v\n", message)
			//校验数据包是否有有效数据
			datalen := message.GetDataLen()
			if datalen == 0 {
				fmt.Printf("数据长度为0，无需读取，msgid:%d\n", message.GetMsgId())
				continue
			}
			//3.第二次读取：读取真实数据
			databuf := make([]byte, datalen)
			cnt, err = io.ReadFull(conn, databuf)
			fmt.Printf("Server========>Client,data:%s,cnt:%d,msgid:%d\n", databuf, cnt, message.GetMsgId())
			//++++++++++++++++++++++++++++++++++++++
			time.Sleep(3 * time.Second)
		}
	}()
	select {}
}
