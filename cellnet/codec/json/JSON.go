package json

import (
	"Chat/cellnet"
	"encoding/json"
)

func init() {
	cellnet.RegisterCodec(jsonCodec{})
}

type jsonCodec struct{}

func (jsonCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (jsonCodec) Decode(data []byte, addressOfVar interface{}) error {
	return json.Unmarshal(data, addressOfVar)
}

func (jsonCodec) Name() string {
	return "json"
}
