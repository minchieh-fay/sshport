package pool

import "sshport/proto"

func (p *Pool) SendSshReset() {
	// 组装一个5字节的byte[]， 类型为MSG_TYPE_SSHRESET
	buf := make([]byte, 5)
	buf[4] = proto.MSG_TYPE_SSHRESET
	p.Write(buf)
}
