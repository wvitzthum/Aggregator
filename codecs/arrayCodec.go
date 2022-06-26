package codecs

import (
	"encoding/json"
	"aggregator/models"
)

type ArrayCodec struct{}

func (c *ArrayCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (c *ArrayCodec) Decode(data []byte) (interface{}, error) {
	var v []models.Event
	err := json.Unmarshal(data, &v)
	return v, err
}
