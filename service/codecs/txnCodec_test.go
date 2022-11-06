package codecs

import (
	"aggregator/service/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecodet(t *testing.T) {
	txn := models.Txn{Op: "xz"}

	b, err := TxnEncodeGOB(&txn)
	assert.NoError(t, err)
	
	txn2, err := TxnDecodeGOB(b)
	assert.NoError(t, err)
	
	assert.Equal(t, txn, txn2)
}
