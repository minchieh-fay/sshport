package proto

import "encoding/binary"

// DataPacket 数据包结构
type DataPacket struct {
	SeqID uint16 // 序列号 (0-65535)
	Data  []byte // 加密后的数据
}

func (dp *DataPacket) Marshal() []byte {
	if len(dp.Data) == 0 {
		return nil
	}
	// 加密
	datatmp, err := Encrypt(dp.Data)
	if err != nil {
		return nil
	}
	dp.Data = datatmp

	// 组装数据包
	dplen := 2 + len(dp.Data)
	data := make([]byte, PacketHeaderSize+dplen)
	//packet.Length = binary.BigEndian.Uint32(headBuf[:4])
	//packet.Type = headBuf[4]
	// 上面2行诗反序列化len和type 现在需要序列化
	binary.BigEndian.PutUint32(data[:4], uint32(dplen))
	data[4] = byte(MSG_TYPE_DATA)
	binary.BigEndian.PutUint16(data[5:7], dp.SeqID)
	copy(data[7:], dp.Data)
	return data
}

func UnmarshalDataPacket(data []byte) (*DataPacket, error) {
	packet := &DataPacket{}
	packet.SeqID = binary.BigEndian.Uint16(data[:2])
	packet.Data = data[2:]
	return packet, nil
}
