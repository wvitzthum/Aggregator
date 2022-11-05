package featureCalculation

import (
	"encoding/json"
	"aggregator/service/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestTxn(val int) models.Txn {
	return models.Txn{
		X: models.Transaction{
			Out: []models.Account{
				{
					Value: val,
				},
			},
		},
	}
}

func TestCalcFeatures(t *testing.T) {
	txns := []models.Txn{getTestTxn(1), getTestTxn(5), getTestTxn(10), getTestTxn(20)}

	fts := CalcFeatures(txns)
	
	assert.Equal(t, fts.Features["TOTAL_SUM"], 36.0, "total sum invalid")
	assert.Equal(t, fts.Features["TOTAL_MEAN"], 9.0, "mean invalid")
	assert.Equal(t, fts.Features["TOTAL_MEDIAN"], 5.0, "median invalid")
}

func TestJsonEncoding(t *testing.T) {
	txns := []models.Txn{getTestTxn(1), getTestTxn(5), getTestTxn(10), getTestTxn(20)}
	fts := CalcFeatures(txns)
	enc, err := json.Marshal(fts)
	assert.NoError(t, err)
	t.Log(string(enc))
}