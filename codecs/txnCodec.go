package codecs

import (
	"encoding/json"

	"aggregator/models"
)

type TxnCodec struct{}

func (c *TxnCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (c *TxnCodec) Decode(data []byte) (interface{}, error) {
	var v models.Txn
	err := json.Unmarshal(data, &v)
	return v, err
}