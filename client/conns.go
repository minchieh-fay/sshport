package client

import (
	"log"
	"net"
	"os"
	"sshport/proto"
	"time"
)

var MAX_CONN_COUNT = 10

func (c *Client) keepconn() {
	for {
		time.Sleep(100 * time.Millisecond)
		if c.pool.GetConnCount() >= MAX_CONN_COUNT {
			time.Sleep(1000 * time.Millisecond)
			continue
		}

		conn, err := net.Dial("tcp", c.Config.ServerAddress)
		if err != nil || conn == nil {
			log.Fatalf("连接服务端失败: %v", err)
			os.Exit(1)
		}
		// 发送auth
		bufs := proto.CreateAuthBufferWithEncrypt()
		if bufs == nil {
			log.Printf("创建auth失败")
			conn.Close()
			continue
		}
		conn.Write(bufs)
		packet, err := proto.DecodePacket(conn)
		if err != nil {
			log.Printf("解码authack失败: %v", err)
			conn.Close()
			continue
		}
		if packet.Type != proto.MSG_TYPE_AUTHACK {
			log.Printf("authack消息类型错误: %d", packet.Type)
			conn.Close()
			continue
		}
		authack := packet.Any.(*proto.AuthAckPacket)
		if !authack.Success {
			log.Printf("authack失败")
			conn.Close()
			continue
		}
		c.pool.AddConn(conn)
	}
}
