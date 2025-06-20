package proto

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"sshport/help"
	"time"
)

var salt = "fmj123456"

type AuthPacket struct {
	Timestamp int64
	Auth      []byte
}

func UnmarshalAuthPacket(data []byte) (*AuthPacket, error) {
	packet := &AuthPacket{}
	packet.Timestamp = int64(binary.BigEndian.Uint64(data[:8]))
	packet.Auth = data[8:]
	return packet, nil
}

func CreateAuthBufferWithEncrypt() []byte {
	timestamp := time.Now().Unix()
	auth := md5.Sum([]byte(fmt.Sprintf("%d%s", timestamp, help.AuthSalt)))
	buf := make([]byte, 8+len(auth))
	binary.BigEndian.PutUint64(buf[:8], uint64(timestamp))
	copy(buf[8:], auth[:])
	buf, err := Encrypt(buf)
	if err != nil {
		return nil
	}
	fullbuf := make([]byte, PacketHeaderSize+len(buf))
	binary.BigEndian.PutUint32(fullbuf[:4], uint32(len(buf)))
	fullbuf[4] = MSG_TYPE_AUTH
	copy(fullbuf[5:], buf)
	return fullbuf
}
