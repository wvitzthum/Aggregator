package featureCalculation

import (
	"aggregator/service/models"
	"log"

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

func CalcFeatures(w []models.Txn) *models.Features {

	features := new(models.Features)

	var err error

	values := make([]float64, len(w)) // the value of each transaction

	for i, txn := range w {
		values[i] = float64(txn.X.Out[0].Value)
	}

	features.SumValue = Sum(values)
	if err != nil {
		log.Fatal(err)
	}

	features.MeanValue = stat.Mean(values, nil)

	features.MedianValue = stat.Quantile(0.5, stat.Empirical, values,nil)

	return features

}