package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//负责解析配置文件

//1. 定义一个配置文件结构
//ip，port，name，version
type Config struct {
	IP            string
	Port          uint32
	Name          string
	Version       string
	WorkSize      int //woker的数量，goroutine
	TaskQueSize   int //每一个消息队列能够容纳请求的最大数量
	ConnAllowSize int //服务器允许最大连接数量
}

//在init函数中加载LoadConfig函数
func init() {
	err := LoadConfig()
	if err != nil {
		fmt.Println("zinx加载配置文件失败:err", err)
		//-1为状态码
		os.Exit(-1)
	}
	fmt.Println("======配置文件信息如下：=====")
	fmt.Printf("%v\n", GlobalConfig)
	fmt.Println("+++++++++++++++++++++++++++++++++=")
}

//2. 加载配置文件//定义一个全局的配置文件结构，用于接收从配置文件中读取的数据
var GlobalConfig Config

func LoadConfig() error {
	fmt.Println("开始读取配置文件...")
	//1.读取配置文件
	//基于server_main.go目录进行寻找
	configInfo, err := ioutil.ReadFile("./conf/conf.json")
	if err != nil {
		return err
	}
	//2. 反序列为Config结构
	//3.配置文件全局唯一，需要定义一个GlobalConfig字段，赋值解析出来的数据
	err = json.Unmarshal(configInfo, &GlobalConfig)
	if err != nil {
		return err
	}
	fmt.Println("读取配置文件")
	return nil
}

//在zinx服务器中使用配置文件
