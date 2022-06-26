package codecs

import (
	"aggregator/models"
	"encoding/json"
	"errors"
)

type eventCodec struct{}

func (e *eventCodec) Encode(value interface{}) ([]byte, error) {
	v, ok := value.(models.Event)
	if !ok {
		return nil, errors.New("Failure encodoing array")
	}
	return json.Marshal(v)
}

func (e *eventCodec) Decode(data []byte) (interface{}, error) {
	var v models.Event
	err := json.Unmarshal(data, &v)
	return v, err
}
