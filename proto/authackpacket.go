package proto

import "encoding/binary"

type AuthAckPacket struct {
	Success bool
}

func UnmarshalAuthAckPacket(data []byte) (*AuthAckPacket, error) {
	packet := &AuthAckPacket{}
	packet.Success = data[0] == 1
	return packet, nil
}

func CreateAuthAckBufferWithEncrypt() []byte {
	buf := make([]byte, 1)
	buf[0] = 1
	buf, err := Encrypt(buf)
	if err != nil {
		return nil
	}
	fullbuf := make([]byte, PacketHeaderSize+len(buf))
	binary.BigEndian.PutUint32(fullbuf[:4], uint32(len(buf)))
	fullbuf[4] = MSG_TYPE_AUTHACK
	copy(fullbuf[5:], buf)
	return fullbuf
}
