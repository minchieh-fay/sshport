package proto

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
)

// 消息类型常量
const (
	MSG_TYPE_DATA     = 1 // 数据消息
	MSG_TYPE_SSHFIN   = 2 // 连接结束消息
	MSG_TYPE_SSHRESET = 3 // 重置seq消息
	MSG_TYPE_AUTH     = 4 // 认证消息
	MSG_TYPE_AUTHACK  = 5 // 认证ACK消息

	MSG_TYPE_MAX = 6 // 标记消息的最大值
)

// 协议常量
const (
	PacketHeaderSize = 5 // 4字节长度 + 1字节类型
	//MaxPacketSize    = 64 * 1024 // 64KB
	AUTH_SECRET = "fmj123" // 认证密钥
)

type Packet struct {
	Length uint32
	Type   uint8
	Any    interface{} // Type对应的一个结构体
}

func DecodePacket(conn net.Conn) (*Packet, error) {
	if conn == nil {
		return nil, errors.New("连接为空")
	}
	packet := &Packet{}
	headBuf := make([]byte, PacketHeaderSize)
	n, err := io.ReadFull(conn, headBuf)
	if err != nil || n != PacketHeaderSize {
		log.Printf("1111")
		return nil, err
	}
	packet.Length = binary.BigEndian.Uint32(headBuf[:4])
	packet.Type = headBuf[4]
	if packet.Type > MSG_TYPE_MAX {
		log.Printf("2222")
		return nil, errors.New("消息类型错误")
	}
	if packet.Length > 0 {
		bodybuf := make([]byte, packet.Length)
		n, err := io.ReadFull(conn, bodybuf)
		if err != nil || n != int(packet.Length) {
			log.Printf("3333")
			return nil, err
		}
		// 解密
		plaintext, err := Decrypt(bodybuf)
		if err != nil {
			log.Printf("4444")
			return nil, err
		}

		// 解密后，根据类型，赋值给对应的结构体
		switch packet.Type {
		case MSG_TYPE_DATA:
			packet.Any, err = UnmarshalDataPacket(plaintext)
			if err != nil {
				log.Printf("5555")
				return nil, err
			}
		case MSG_TYPE_SSHFIN: // 不用处理
		case MSG_TYPE_SSHRESET: // 不用处理
		case MSG_TYPE_AUTH:
			packet.Any, err = UnmarshalAuthPacket(plaintext)
			if err != nil {
				log.Printf("6666")
				return nil, err
			}
		case MSG_TYPE_AUTHACK:
			packet.Any, err = UnmarshalAuthAckPacket(plaintext)
			if err != nil {
				log.Printf("7777")
				return nil, err
			}
		default:

			return nil, errors.New("消息类型错误")
		}
	}
	return packet, nil
}
