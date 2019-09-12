package net

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

//1. 每个测试文件需要以_test.go结尾
//2. 每个测试文件需要引用testing包
//3. 每个测试的函数需要以Test开头。
func TestDataPack(t *testing.T) {
	fmt.Printf("TestDataPack called!")
	//测试data pack 与unpack函数
	//	server
	go func() {
		//	监听
		listener, err := net.Listen("tcp", "0.0.0.0:8888")
		if err != nil {
			t.Errorf("net.Listen err :%v", err)
			return
		}
		conn, err := listener.Accept()
		if err != nil {
			t.Errorf("listener.Accept err:%v\n", err)
			return
		}
		//开始读取
		for {
			//1. 第一次读取：读取数据头长度的数据（8字节）2个uint32
			headBuffer := make([]byte, 8)
			//2. 拆包, 把数据的长度(N)，与消息的类型解析出来
			//c.Read()
			//a. 只能一次性读取，没有指定长度, 如果有网络延迟，可能没有读取我们需要的长度的数据就返回
			//b. 解决办法，使用io.ReadFull来读取数据，这个函数可以读取指定的buf长度的数据，如果未读取完毕，则不返回
			cnt, err := io.ReadFull(conn, headBuffer)
			if err != nil {
				t.Errorf("io.ReadFull err；%v\n", err)
				return
			}
			fmt.Printf("读取数据头的长度：%d\n", cnt)
			//	拆包
			dp := NewDataPack()
			err, message := dp.Unpack(headBuffer)
			fmt.Printf("拆包之后的数据详情：%v\n", message)
			//校验数据包是否有有效数据
			dataLen := message.GetDataLen()
			if dataLen == 0 {
				fmt.Printf("数据长度为0，无需读取，msgid:%d\n", message.GetMsgId())
				continue
			}
			//	用于存储消息体即真实的数据
			databuf := make([]byte, dataLen)
			//第二次读取：读取真实数据长度（长度为N)
			cnt, err = io.ReadFull(conn, databuf)
			fmt.Printf("Server《========Client,data:%s,cnt:%d,msgid:%d\n", databuf, cnt, message.GetMsgId())
		}

	}()
	//client
	go func() {
		//封包，发送
		//1.准备数据（封包）
		data0 := []byte{}
		data1 := []byte("你好")
		data2 := []byte("hello world")
		data3 := []byte("中秋即将到来")
		//2.创建message
		msg0 := NewMessage(data0, uint32(len(data0)), 0)
		msg1 := NewMessage(data1, uint32(len(data1)), 1)
		msg2 := NewMessage(data2, uint32(len(data2)), 2)
		msg3 := NewMessage(data3, uint32(len(data3)), 3)
		//3.对message进行封包
		dp := NewDataPack()
		_, infobytes0 := dp.Pack(msg0)
		_, infobytes1 := dp.Pack(msg1)
		_, infobytes2 := dp.Pack(msg2)
		_, infobytes3 := dp.Pack(msg3)
		//4.将消息的字节流拼接到一起，一次性发送给服务器//把多个包黏在一起，一起发送
		infosend := append(infobytes0, infobytes1...)
		infosend = append(infosend, infobytes2...)
		infosend = append(infosend, infobytes3...)
		//	5.发送
		conn, err := net.Dial("tcp", "127.0.0.1:8888")
		if err != nil {
			t.Errorf("client dial err:%v\n", err)
			return
		}
		cnt, err := conn.Write(infosend)
		if err != nil {
			t.Errorf("client send err:%v\n", err)
			return
		}
		fmt.Println("Client====>Server cnt:", cnt)

	}()
	time.Sleep(2 * time.Second)
}
