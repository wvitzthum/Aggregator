package models

import (
	"github.com/lovoo/goka"
)

type Topic struct {
	Stream  goka.Stream
	Codec goka.Codec
}

type Txn struct {
	Op string `json:"op"`
	X  struct {
		Hash   string `json:"hash"`
		Inputs []struct {
			PrevOut struct {
				Addr    string `json:"addr"`
				N       int    `json:"n"`
				Script  string `json:"script"`
				Spent   bool   `json:"spent"`
				TxIndex int    `json:"tx_index"`
				Type    int    `json:"type"`
				Value   int    `json:"value"`
			} `json:"prev_out"`
			Script   string `json:"script"`
			Sequence int64  `json:"sequence"`
		} `json:"inputs"`
		LockTime int `json:"lock_time"`
		Out      []struct {
			Addr    string `json:"addr"`
			N       int    `json:"n"`
			Script  string `json:"script"`
			Spent   bool   `json:"spent"`
			TxIndex int    `json:"tx_index"`
			Type    int    `json:"type"`
			Value   int    `json:"value"`
		} `json:"out"`
		RelayedBy string `json:"relayed_by"`
		Size      int    `json:"size"`
		Time      int64   `json:"time"`
		TxIndex   int    `json:"tx_index"`
		Ver       int    `json:"ver"`
		VinSz     int    `json:"vin_sz"`
		VoutSz    int    `json:"vout_sz"`
	} `json:"x"`
}
