package codecs

import (
	"encoding/json"
)

type ArrayCodec struct{}

func (c *ArrayCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (c *ArrayCodec) Decode(data []byte) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	return v, err
}
