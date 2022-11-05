package codecs

import (
	"encoding/json"
	"aggregator/service/models"
)
type FeaturesCodec struct{}

func (c *FeaturesCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (c *FeaturesCodec) Decode(data []byte) (interface{}, error) {
	var v models.Features
	err := json.Unmarshal(data, &v)
	return v, err
}