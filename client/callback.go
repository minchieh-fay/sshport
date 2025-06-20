package client

import "sshport/help"

func (c *Client) callback(t int) {
	switch t {
	case help.CALLBACK_TYPE_SSHFIN:
		c.sshconn.Close()
		c.sshconn = nil
	}
}
