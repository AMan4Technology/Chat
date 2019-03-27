package packet

import (
	"bytes"
	"encoding/binary"
)

type Reader struct {
	raw    []byte
	reader *bytes.Reader
}

func (r *Reader) Init(body []byte) {
	r.raw = body
	r.reader = bytes.NewReader(body)
}

func (r Reader) Raw() []byte {
	return r.raw
}

func (r Reader) RemainLen() int {
	return r.reader.Len()
}

func (r Reader) RemainBytes() []byte {
	return r.raw[len(r.raw)-r.reader.Len():]
}

func (r Reader) ReadValue(addressOfV interface{}) error {
	return binary.Read(r.reader, binary.LittleEndian, addressOfV)
}

func (r Reader) ReadString(addressOfV *string) error {
	var size uint16
	if err := r.ReadValue(&size); err != nil {
		return err
	}
	body := make([]byte, size)
	if err := r.ReadValue(&body); err != nil {
		return err
	}
	*addressOfV = string(body)
	return nil
}
