package cellnet

import (
	"bytes"
	"fmt"
	"path"
	"reflect"
)

var (
	metaWithName = make(map[string]*MessageMeta)
	metaWithID   = make(map[uint32]*MessageMeta)
	metaWithType = make(map[reflect.Type]*MessageMeta)
)

func RegisterMessageMeta(codecName string, name string, msgType reflect.Type, id uint32) {
	meta := &MessageMeta{
		Type:  msgType,
		Name:  name,
		ID:    id,
		Codec: GetCodec(codecName),
	}
	if meta.Codec == nil {
		panic("codec not register! " + codecName)
	}
	if MessageMetaByName(name) != nil {
		panic("duplicate message meta register by name: " + name)
	}
	if MessageMetaByID(id) != nil {
		panic(fmt.Sprintf("duplicate message meta register by id: %d", id))
	}
	if MessageMetaByType(msgType) != nil {
		panic(fmt.Sprintf("duplicate message meta register by type: %s", msgType.Name()))
	}
	metaWithName[name] = meta
	metaWithID[id] = meta
	metaWithType[msgType] = meta
}

func MessageMetaByName(name string) *MessageMeta {
	if meta, ok := metaWithName[name]; ok {
		return meta
	}
	return nil
}

func MessageMetaByID(id uint32) *MessageMeta {
	if meta, ok := metaWithID[id]; ok {
		return meta
	}
	return nil
}

func MessageMetaByType(msgType reflect.Type) *MessageMeta {
	if meta, ok := metaWithType[msgType]; ok {
		return meta
	}
	return nil
}

func MessageFullName(msgType reflect.Type) string {
	if msgType == nil {
		panic("empty msg type")
	}
	if msgType.Kind() == reflect.Ptr {
		msgType = msgType.Elem()
	}
	var buffer bytes.Buffer
	buffer.WriteString(path.Base(msgType.PkgPath()))
	buffer.WriteByte('.')
	buffer.WriteString(msgType.Name())
	return buffer.String()
}

type MessageMeta struct {
	Type  reflect.Type
	Name  string
	ID    uint32
	Codec Codec
}
