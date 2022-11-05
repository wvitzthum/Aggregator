package featureCalculation

import (
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
	
	assert.Equal(t, fts.SumValue, 36.0)
	assert.Equal(t, fts.MeanValue, 9.0)
	assert.Equal(t, fts.MedianValue, 5.0)
}