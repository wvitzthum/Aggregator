package views

import (
	"aggregator/service/codecs"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lovoo/goka"
)

func RunViews(brokers []string) {
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
	root.HandleFunc("/window/{key}", getByID(view))
	root.HandleFunc("/window/", getAll(view))

	fmt.Println("View opened at http://localhost:9095/")
	go http.ListenAndServe(":9095", root)

	err = view.Run(context.Background())
	log.Println(err)
}
