package client

import (
	"log"
	"time"
)

func (c *Client) readSsh() {
	for {
		conn := c.sshconn
		if conn == nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		buf := make([]byte, 8192)
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("读取ssh连接失败: %v", err)
			conn.Close()
			c.sshconn = nil
			c.pool.SendSshFin()
			continue
		}
		c.pool.Write(buf[:n])
	}
}
