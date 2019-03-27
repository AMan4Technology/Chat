package packet

import (
	"bytes"
	"encoding/binary"
)

type Writer struct {
	buffer bytes.Buffer
}

func (w Writer) Len() uint16 {
	return uint16(w.buffer.Len())
}

func (w *Writer) WriteValue(value interface{}) error {
	return binary.Write(&w.buffer, binary.LittleEndian, value)
}

func (w *Writer) Raw() []byte {
	return w.buffer.Bytes()
}

func (w *Writer) WriteString(value string) error {
	if err := w.WriteValue(len(value)); err != nil {
		return err
	}
	if err := w.WriteValue(value); err != nil {
		return err
	}
	return nil
}
