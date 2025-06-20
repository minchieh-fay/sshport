package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"sshport/client"
	"sshport/help"
	"sshport/server"
	"time"
)

var Config = &help.ConfigInfo{
	// 服务端参数
	SshAddress: "127.0.0.1:22",

	// 客户端参数
	ServerAddress: "",

	// 通用参数
	Port:  0,
	Key:   "",
	Help:  false,
	Debug: false,
}

func init() {
	flag.IntVar(&Config.Port, "d", 0, "端口")
	flag.StringVar(&Config.ServerAddress, "h", "", "服务端地址ip:port")
	flag.StringVar(&Config.Key, "k", "awefeawgaw", "加密密钥")
	flag.StringVar(&Config.SshAddress, "sh", "127.0.0.1:22", "本地需要代理的ssh地址")
	flag.BoolVar(&Config.Help, "help", false, "帮助")
	flag.Parse()

	if Config.Help {
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func test() {
	Config.ServerAddress = "127.0.0.1:5566"
	Config.Port = 5566
	Config.SshAddress = "10.35.148.167:22"
	Config.Key = "awefeawgaw"
	Config.Help = false
	Config.Debug = true
	go server.NewServer(Config).Start()
	time.Sleep(200 * time.Millisecond)
	client.NewClient(Config).Start()
}

func main() {
	// 设置单线程
	runtime.GOMAXPROCS(1)
	if 1 == 1 {
		test()
		return
	}

	log.Printf("Config: %v", Config)
	if Config.ServerAddress != "" {
		client.NewClient(Config).Start()
	} else {
		server.NewServer(Config).Start()
	}
}
