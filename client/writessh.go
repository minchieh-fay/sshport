package client

import (
	"time"
)

func (c *Client) writeSsh() {
	for {
		d := c.pool.Read()
		if d == nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		if c.sshconn == nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		_, err := c.sshconn.Write(d)
		if err != nil {
			c.sshconn.Close()
			c.sshconn = nil
			continue
		}
	}
}
