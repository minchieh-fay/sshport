package pool

import (
	"sshport/proto"
	"time"
)

func (p *Pool) Read() []byte {
	for {
		dp := <-p.ch
		if dp == nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		return dp.Data
	}
}

func (p *Pool) Write(data []byte) {
	if len(p.conns) == 0 {
		return
	}
	dp := &proto.DataPacket{
		SeqID: p.sendseq,
		Data:  data,
	}
	p.sendseq++
	for {
		for _, c := range p.conns {
			if c.ref > 0 {
				continue
			}
			c.ref = 1
			go func() {
				_, err := c.conn.Write(dp.Marshal())
				if err != nil {
					// 删除conn
					for i, c := range p.conns {
						if c.conn == c.conn {
							p.conns = append(p.conns[:i], p.conns[i+1:]...)
							break
						}
					}
				}
				c.ref = 0
			}()
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}
