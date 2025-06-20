package server

import (
	"crypto/md5"
	"fmt"
	"log"
	"net"
	"sshport/help"
	"sshport/proto"
	"time"
)

func (s *Server) authConn(conn net.Conn) bool {
	packet, err := proto.DecodePacket(conn)
	if err != nil {
		log.Printf("auth解码消息失败: %v", err)
		return false
	}
	if packet == nil {
		log.Printf("auth消息为空")
		return false
	}
	if packet.Type != proto.MSG_TYPE_AUTH {
		log.Printf("auth消息类型错误: %d", packet.Type)
		return false
	}
	timestamp := packet.Any.(*proto.AuthPacket).Timestamp
	auth := packet.Any.(*proto.AuthPacket).Auth

	// 检查时间戳是否过期（10秒内有效）
	now := time.Now().Unix()
	if abs(now-timestamp) > 10 {
		log.Printf("认证时间戳过期: 当前时间=%d, 包时间=%d", now, timestamp)
		return false
	}

	// md5(timestamp+slat )
	auth2 := md5.Sum([]byte(fmt.Sprintf("%d%s", timestamp, help.AuthSalt)))

	// 检查认证密钥
	if string(auth) != string(auth2[:]) {
		log.Printf("认证密钥错误: %s", auth)
		return false
	}
	// 发送authack
	bufs := proto.CreateAuthAckBufferWithEncrypt()
	if bufs == nil {
		log.Printf("创建authack失败")
		return false
	}
	conn.Write(bufs)

	return true
}

// abs 计算绝对值
func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
