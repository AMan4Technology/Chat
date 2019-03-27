package proto

import (
	"Chat/cellnet"
	_ "Chat/cellnet/codec/json"
	"reflect"
)

func init() {
	cellnet.RegisterMessageMeta("json", "proto.ChatREQ", reflect.TypeOf((*ChatREQ)(nil)).Elem(), 1)
	cellnet.RegisterMessageMeta("json", "proto.ChatACK", reflect.TypeOf((*ChatACK)(nil)).Elem(), 2)
}

type ChatREQ struct {
	Content string
}

type ChatACK struct {
	Content string
	ID      int64
}
