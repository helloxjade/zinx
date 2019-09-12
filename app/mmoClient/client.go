package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	pb2 "github.com/helloxjade/zinx/app/mmoClient/pb"
	"io"
	"math/rand"
	"net"
	"time"
)

type Client struct {
	//	对应Service 中的player
	//属性
	//1.唯一的id pid
	Pid int
	//	2. 原生的Conn,不是zinx,因为这和zinx 框架 无关
	Conn net.Conn
	//	3.玩家位置
	X      int
	Y      int       //高度
	Z      int       //纵轴
	V      int       //面部朝向
	online chan bool //标识客户端线
}

//创建一个Client ，pid 和position 在client上线后，会由服务器主动发送过来
func NewClient(ip string, port int) *Client {
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}
	return &Client{
		Pid:    0,
		Conn:   conn,
		X:      0,
		Y:      0,
		Z:      0,
		V:      0,
		online: make(chan bool),
	}
}
func (client *Client) Start() {
	//接收服务器数据
	go func() {
		for {
			fmt.Println("客户端处理业务...")
			//	读取服务器发来的数据===》SyncPid
			//1.读取8字节
			//第一次读取，读取数据头
			headData := make([]byte, 8)
			if _, err := io.ReadFull(client.Conn, headData); err != nil {
				fmt.Println("读取数据头8字节失败：", err)
				return
			}
			//拆包
			dp := NewDataPack()
			msg, err := dp.UnPack(headData)
			if err != nil {
				fmt.Println("dp Unpack err:", err)
				return
			}
			//	第一次读数据长度
			dataLen := msg.len
			if dataLen == 0 {
				fmt.Println("无具体数据，不用再读取！！")
				continue
			}
			//第二次读取：读取真实长度
			realData := make([]byte, dataLen)
			if _, err := io.ReadFull(client.Conn, realData); err != nil {
				fmt.Println("读取真实数据失败：", err)
				return
			}
			//将读取的数据拼接到Message 结构中
			msg.data = realData
			//调用处理具体消息的业务
			client.HandleMsg(msg)
		}
	}()
	//读取到用户上线后，服务器才开始发送请求：聊天，移动
	for {
		select {
		case <-client.online:
			go func() {
				for {
					//向服务器发送请求：聊天，移动
					client.robotAction() //新建一个函数处理聊天、移动
					time.Sleep(1 * time.Second)
				}
			}()
		}
	}
	//不断的发请求
	select {}
}

//处理发送的消息
func (client *Client) HandleMsg(message *Message) {
	fmt.Println("得到Message信息，msgid:", message.msgid)
	//	同步pid
	if message.msgid == 1 {
		fmt.Println("获取玩家pid逻辑，msgid:", message.msgid)
		//	解析出proto内容，得到pid,赋值给client
		var syncPid pb2.SyncPid
		err := proto.Unmarshal(message.data, &syncPid)
		if err != nil {
			fmt.Println("proto 解码失败：", err)
			return
		}
		fmt.Println("获取player id:", syncPid.Pid)
		client.Pid = int(syncPid.Pid)
	} else if message.msgid == 200 {
		fmt.Println("获取广播逻辑，msgid:", message.msgid)
		// 广播消息
		var broadcastData pb2.BroadCast
		err := proto.Unmarshal(message.data, &broadcastData)
		if err != nil {
			fmt.Println("proto解码失败：", err)
			return
		}
		//判断具体的业务类型：1-聊天  2-位置 4-玩家移动
		if broadcastData.Tp == 2 && broadcastData.Pid == int32(client.Pid) {
			//	服务给自己分配了一个位置信息，更新坐标。
			client.X = int(broadcastData.GetP().X)
			client.Y = int(broadcastData.GetP().Y)
			client.Z = int(broadcastData.GetP().Z)
			client.V = int(broadcastData.GetP().V)
			fmt.Printf("玩家 id: %d 已经成功上线，坐标：X :%d Y:%d Z:%d V:%d\n", client.Pid, client.X, client.Y, client.Z, client.V)

			//	服务器发送完pid,发送完玩家位置之后，表明玩家上线成功
			client.online <- true
		} else if broadcastData.Tp == 1 && broadcastData.Pid == int32(client.Pid) {
			fmt.Println("世界聊天：玩家：%d 说的话是：%s", client.Pid, broadcastData.GetContent())
		}
	}
}

//模拟客户端的随机请求
func (client *Client) robotAction() {
	//	提供一个随机数，可以随机得到两个值
	randNum := rand.Intn(2) //0,1
	if randNum == 0 {
		//	聊天 ：msgid 2
		talkInfo := fmt.Sprintf("大家好，我是玩家%d")
		proto_talk := pb2.Talk{
			Content: talkInfo,
		}
		//编码病发送
		client.SendMsg(2, &proto_talk)

	} else if randNum == 1 {
		//	msgid=3 移动
		// 自由移动
		x := client.X
		z := client.Z
		randPos := rand.Intn(2)
		if randPos == 0 {
			//	0，x,z加上一个数据
			x += rand.Intn(10)
			z += rand.Intn(10)
		} else {
			// 1 , z减去一个数据
			x -= rand.Intn(10)
			z -= rand.Intn(10)
		}
		//纠正坐标
		if x > 410 {
			x = 410
		} else if x < 85 {
			x = 85
		}
		if z > 400 {
			z = 400
		} else if z < 75 {
			z = 75
		}
		//面朝方向
		randv := rand.Intn(2)
		v := client.V //下面的数是随机的
		if randv == 0 {
			v = 25
		} else {
			v = 350
		}

		//将最新位置打包成proto结构发给服务器
		proto_position := pb2.Position{
			X: float32(x),
			Y: float32(client.Y),
			Z: float32(z),
			V: float32(v),
		}
		fmt.Printf("2222222=====>玩家：%d的新位置为：x:%d,y:%d,z:%d,v:%v\n", client.Pid, x, client.Y, z, v)
		client.SendMsg(3, &proto_position)
	}
}

func (client *Client) SendMsg(msgid uint32, data proto.Message) {
	//	向服务器发送数据
	//1.proto结构数据编码成二进制流
	binaryInfo, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("SendMsg proto Marshal err:", err)
		return
	}
	//2.封装成pack
	dp := NewDataPack()
	sendIndo, err := dp.Pack(NewMessage(msgid, binaryInfo))
	if err != nil {
		fmt.Println("SendMsg Pack err:", err)
		return
	}

	//	发送
	cnt, err := client.Conn.Write(sendIndo)
	if err != nil {
		fmt.Println("SendMsg Conn Write err:", err)
		return
	}
	fmt.Println("Client====>Server,cnt", cnt)
}
