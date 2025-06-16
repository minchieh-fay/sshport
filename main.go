package main

import (
	"flag"
	"log"
	"os"
	"time"
)

// 参数读取

// 服务端参数
// 1. -d 服务端口
// 2. -sh 本地需要代理的ssh地址  默认127.0.0.1:22

// 客户端参数
// 1. -h 服务端地址ip:port

// 通用参数
// 1. 加密密钥  // 流式传输的时候 需要加密 加在数据上面的 和quic本身的加密无关

// 如果参数含有-d就是服务端
// 如果参数含有-h就是客户端

type ConfigInfo struct {
	Port          int // 如果是服务端  那么，他就是quic的端口， 如果是客户端那么就是tcp端口
	ServerAddress string
	Key           string
	Help          bool
	SshAddress    string
}

var Config = &ConfigInfo{
	// 服务端参数
	SshAddress: "127.0.0.1:22",

	// 客户端参数
	ServerAddress: "",

	// 通用参数
	Port: 0,
	Key:  "",
	Help: false,
}

func init() {
	flag.IntVar(&Config.Port, "d", 0, "端口")
	flag.StringVar(&Config.ServerAddress, "h", "", "服务端地址ip:port")
	flag.StringVar(&Config.Key, "k", "", "加密密钥")
	flag.StringVar(&Config.SshAddress, "sh", "127.0.0.1:22", "本地需要代理的ssh地址")
	flag.BoolVar(&Config.Help, "help", false, "帮助")
	flag.Parse()

	if Config.Help {
		flag.PrintDefaults()
		os.Exit(0)
	}
}

var debug = false

func test() {
	debug = true
	Config.Port = 5566
	Config.ServerAddress = "127.0.0.1:5566"
	Config.SshAddress = "10.35.148.167:22"
	Config.Key = "1234567890"

	server := NewQuicServer(Config)
	go server.Start()
	client := NewQuicClient(Config)
	client.Start()
	for {
		time.Sleep(time.Second * 1)
	}
}

func main() {
	log.Println("start")

	// if 1 == 1 {
	// 	test()
	// }
	if Config.ServerAddress == "" { // 进入服务模式
		log.Println("进入服务模式")
		server := NewQuicServer(Config)
		server.Start()
	} else { // 进入客户端模式
		log.Println("进入客户端模式")
		client := NewQuicClient(Config)
		client.Start()
	}
}
