package packet

import (
	"encoding/binary"
	"io"
)

const LengthSize = 2

func RecvVariableLengthPacket(reader io.Reader) (pktReader Reader, err error) {
	sizeBuffer := make([]byte, LengthSize)
	_, err = io.ReadFull(reader, sizeBuffer)
	if err != nil {
		return
	}
	size := binary.LittleEndian.Uint16(sizeBuffer)
	body := make([]byte, size)
	_, err = io.ReadFull(reader, body)
	pktReader.Init(body)
	return
}

func SendVariableLengthPacket(writer io.Writer, pktWriter Writer) error {
	buffer := make([]byte, LengthSize+pktWriter.Len())
	binary.LittleEndian.PutUint16(buffer, pktWriter.Len())
	copy(buffer[LengthSize:], pktWriter.Raw())
	if _, err := writer.Write(buffer); err != nil {
		return err
	}
	return nil
}
