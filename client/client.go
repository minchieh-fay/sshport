package client

import (
	"fmt"
	"log"
	"net"
	"sshport/help"
	"sshport/pool"
	"time"
)

type Client struct {
	Config  *help.ConfigInfo
	pool    *pool.Pool
	sshconn net.Conn
}

func NewClient(config *help.ConfigInfo) *Client {
	return &Client{
		Config: config,
		pool:   pool.NewPool(config, "client"),
	}
}

func (c *Client) Start() {
	go c.keepconn()
	c.pool.SetCallback(c.callback)

	// 刚启动， 发送MSG_TYPE_SSHRESET让服务的将seq置0
	c.pool.SendSshReset()

	time.Sleep(1 * time.Second)

	go c.writeSsh()
	go c.readSsh()

	port := c.Config.Port
	if c.Config.Debug {
		port = port + 1
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("监听端口失败: %v", err)
	}
	defer listener.Close()

	log.Printf("客户端已启动，监听端口: %d", c.Config.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("接受连接失败: %v", err)
			continue
		}
		if c.sshconn != nil {
			c.sshconn.Close()
		}
		c.sshconn = conn
	}
}
