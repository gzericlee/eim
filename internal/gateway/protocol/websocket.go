package protocol

import (
	"encoding/binary"
)

type websocketCodec struct {
}

func (its *websocketCodec) Encode(cmd int, buf []byte) []byte {
	msgLen := 2
	if buf != nil {
		msgLen += len(buf)
	}
	data := make([]byte, msgLen)
	binary.BigEndian.PutUint16(data[0:2], uint16(cmd))
	if buf != nil {
		copy(data[2:msgLen], buf)
	}
	return data
}

func (its *websocketCodec) Decode(data []byte) (int, []byte) {
	cmd := binary.BigEndian.Uint16(data[0:2])
	return int(cmd), data[2:len(data)]
}
