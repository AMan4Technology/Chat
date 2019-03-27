package cellnet

import (
	"errors"
	"reflect"
)

var codecWithName = make(map[string]Codec)

func RegisterCodec(c Codec) {
	if _, ok := codecWithName[c.Name()]; ok {
		panic("duplicate codec: " + c.Name())
	}
	codecWithName[c.Name()] = c
}

func GetCodec(name string) Codec {
	return codecWithName[name]
}

var (
	MsgNotFoundErr   = errors.New("msg not found")
	CodecNotFoundErr = errors.New("codec not found")
)

func EncodeMessage(msg interface{}) (data []byte, msgID uint32, err error) {
	meta := MessageMetaByType(reflect.TypeOf(msg))
	if meta == nil {
		return nil, 0, MsgNotFoundErr
	}
	msgID = meta.ID
	data, err = meta.Codec.Encode(msg)
	return
}

func DecodeMessage(msgID uint32, data []byte) (interface{}, error) {
	meta := MessageMetaByID(msgID)
	if meta == nil {
		return nil, MsgNotFoundErr
	}
	addressOfMsg := reflect.New(meta.Type).Interface()
	err := meta.Codec.Decode(data, addressOfMsg)
	if err != nil {
		return nil, err
	}
	return addressOfMsg, nil
}

type Codec interface {
	Encode(value interface{}) ([]byte, error)
	Decode(data []byte, addressOfVar interface{}) error
	Name() string
}
