package models

import (
	"github.com/lovoo/goka"
)

type Features struct {
	SumValue    float64
	MeanValue   float64
	MedianValue float64
	Degree      float64
}

type Topic struct {
	Stream  *goka.Stream
	Codec 	goka.Codec
}

type Account struct {
	Addr    string `json:"addr"`
	N       int    `json:"n"`
	Script  string `json:"script"`
	Spent   bool   `json:"spent"`
	TxIndex int    `json:"tx_index"`
	Type    int    `json:"type"`
	Value   int    `json:"value"`
}

type Input struct {
	PrevOut Account `json:"prev_out"`
	Script   string `json:"script"`
	Sequence int64  `json:"sequence"`
}

type Transaction struct {
	Hash   		string	 	`json:"hash"`
	Inputs 		[]Input  	`json:"inputs"`
	LockTime 	int		 	`json:"lock_time"`
	Out      	[]Account	`json:"out"`
	RelayedBy	string 		`json:"relayed_by"`
	Size     	int   		`json:"size"`
	Time     	int64 		`json:"time"`
	TxIndex  	int   		`json:"tx_index"`
	Ver      	int   		`json:"ver"`
	VinSz    	int   		`json:"vin_sz"`
	VoutSz   	int   		`json:"vout_sz"`
}	

type Txn struct {
	Op string `json:"op"`
	X  Transaction `json:"x"`
}
