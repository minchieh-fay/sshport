package pool

import (
	"log"
	"net"
	"sshport/help"
	"sshport/proto"
)

type Conn struct {
	conn net.Conn
	ref  int
}

type Pool struct {
	Config   *help.ConfigInfo
	conns    []*Conn
	pkts     map[uint16]*proto.DataPacket
	seq      uint16
	sendseq  uint16
	ch       chan *proto.DataPacket
	callback func(t int)
	strtype  string
}

func NewPool(config *help.ConfigInfo, strtype string) *Pool {
	return &Pool{
		Config:  config,
		pkts:    make(map[uint16]*proto.DataPacket),
		seq:     0,
		ch:      make(chan *proto.DataPacket, 1000),
		strtype: strtype,
	}
}

func (p *Pool) Start() {

}

func (p *Pool) SendSshFin() {

}

func (p *Pool) GetConnCount() int {
	return len(p.conns)
}

func (p *Pool) AddConn(conn net.Conn) {
	p.conns = append(p.conns, &Conn{
		conn: conn,
		ref:  0,
	})
	go p.handleConn(conn)
}

func (p *Pool) handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
		for i, c := range p.conns {
			if c.conn == conn {
				p.conns = append(p.conns[:i], p.conns[i+1:]...)
				break
			}
		}
	}()

	for {
		packet, err := proto.DecodePacket(conn)
		if err != nil {
			log.Printf("解码消息失败: %v", err)
			return
		}
		switch packet.Type {
		case proto.MSG_TYPE_DATA:
			pkt := packet.Any.(*proto.DataPacket)
			p.AddPkt(pkt)
		case proto.MSG_TYPE_SSHFIN:
			p.callback(help.CALLBACK_TYPE_SSHFIN)
		case proto.MSG_TYPE_SSHRESET:
			p.seq = 0
			p.sendseq = 0
		}
	}
}

func (p *Pool) AddPkt(dp *proto.DataPacket) {
	if dp == nil {
		return
	}
	p.pkts[dp.SeqID] = dp
	for {
		if _, ok := p.pkts[p.seq]; ok {
			p.ch <- p.pkts[p.seq]
			delete(p.pkts, p.seq)
			p.seq++
		} else {
			break
		}
	}
}

func (p *Pool) SetCallback(callback func(t int)) {
	p.callback = callback
}
