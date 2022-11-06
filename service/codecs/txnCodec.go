package codecs

import (
	"bytes"
	"encoding/gob"
	"encoding/json"

	"aggregator/service/models"
)


type WindowCodec struct{}

func (c *WindowCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (c *WindowCodec) Decode(data []byte) (interface{}, error) {
	var v string
	err := json.Unmarshal(data, &v)
	return v, err
}


type TxnCodec struct{}

func (c *TxnCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (c *TxnCodec) Decode(data []byte) (interface{}, error) {
	var v models.Txn
	err := json.Unmarshal(data, &v)
	return v, err
}

func TxnEncodeGOB(value *models.Txn) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(value)
	return buf.Bytes(), err
}

func TxnDecodeGOB(data []byte) (*models.Txn, error) {
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	var v models.Txn
	err := dec.Decode(&v)
	return &v, err
}
