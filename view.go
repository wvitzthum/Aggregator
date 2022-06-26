package main

import (
	"aggregator/codecs"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lovoo/goka"
)

func runView() {
	log.Println("running view")

	view, err := goka.NewView(brokers,
		goka.GroupTable("window"),
		new(codecs.ArrayCodec),
	)
	if err != nil {
		panic(err)
	}

	root := mux.NewRouter()

	// /{key} returns the full history for that key ordered by descending
	// time
	root.HandleFunc("/{key}", func(w http.ResponseWriter, r *http.Request) {

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

		// write the serialised window to the responseWriter
		w.Write(data)
	})

	fmt.Println("View opened at http://localhost:9095/")
	go http.ListenAndServe(":9095", root)

	err = view.Run(context.Background())
	log.Println(err)
}
