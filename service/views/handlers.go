package views

import (
	"aggregator/service/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lovoo/goka"
)



func getByID(view *goka.View) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the stored window
		value, err := view.Get(mux.Vars(r)["key"])
		if err != nil {
			log.Fatal(err)
		}

		if value == nil {
			w.Write([]byte("value is nil"))
			return
		}

		data, err := json.Marshal(value)
		if err != nil {
			log.Fatal(err)
		}

		// write the serialized window to the responseWriter
		w.Write(data)
	}
}

func getAll(view *goka.View) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the stored window
		iter, err := view.Iterator()
		if err != nil {
			w.Write([]byte("Iterator error: " + err.Error()))
			return
		}

		values := map[string][]models.Txn{}
		for iter.Next(){
			v, err := iter.Value()
			if v == nil {
				// TODO: add logging here
				continue
			} else if err != nil {
				// TODO: add logging here
				continue
			}
			// TODO: add check if values in map are not overwritten
			iterVals := v.([]models.Txn) 
			values[iterVals[0].X.Inputs[0].PrevOut.Addr] = iterVals
		}
		data, err := json.Marshal(values)
		if err != nil {
			log.Fatal(err)
		}

		// write the serialized window to the responseWriter
		w.Write(data)
	}
}