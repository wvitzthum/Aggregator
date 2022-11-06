package featureCalculation

import (
	"aggregator/service/models"
	"log"
	"math"

	"gonum.org/v1/gonum/stat"
)


func Values(t []models.Txn) []float64 {
	res := make([]float64, len(t))
	for _, v := range(t) {
		for _, v2 := range(v.X.Out) {
			res = append(res, float64(v2.Value))
		}
	}
	return res
}

func Sum(w []float64) float64 {
	res := 0.0
	for _, v := range(w) {
		res += v
	}
	return res
}

func Variance(v []float64, w []float64) float64{
	res := stat.Variance(v, w)
	if math.IsNaN(res) {
		return 0.0
	}
	return res
}

func UpdateRelsAndTxns(txns []models.Txn, fts *models.Features) {
	for _, v := range(txns) {
		fts.Transactions = append(fts.Transactions, v.X.Hash)
		for _, v2 := range(v.X.Inputs) {
			fts.Relationships.In = append(fts.Relationships.In, v2.PrevOut.Addr)
		}
		for _, v2 := range(v.X.Out) {
			fts.Relationships.Out = append(fts.Relationships.Out, v2.Addr)
		}
	}
	
}

func CalcFeatures(w []models.Txn) *models.Features {

	features := &models.Features{Features: map[string]float64{}}

	var err error

	values := make([]float64, len(w)) // the value of each transaction

	for i, txn := range w {
		values[i] = float64(txn.X.Out[0].Value)
	}

	features.Features["TOTAL_SUM"] = Sum(values)
	if err != nil {
		log.Fatal(err)
	}

	features.Features["TOTAL_MEAN"] = stat.Mean(values, nil)
	features.Features["TOTAL_MEDIAN"] = stat.Quantile(0.5, stat.Empirical, values,nil)
	features.Features["TOTAL_VARIANCE"] = Variance(values, nil)
	features.Features["TOTAL_COUNT"] = float64(len(w))
	UpdateRelsAndTxns(w, features)
	return features

}