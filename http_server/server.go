package http_server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/tie304/investment/database"
)

const port = ":8000"

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	// get all assets
	var assets []database.Asset
	database.Database.Model(&assets).Select()
	data, _ := json.Marshal(assets)
	w.Write(data)
}

func assetHandler(w http.ResponseWriter, r *http.Request) {
	// create new asset
	if r.Method == http.MethodPost {
		var asset database.Asset
		err := json.NewDecoder(r.Body).Decode(&asset)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Malformed asset data"))
		}
		_, err = database.Database.Model(&asset).Insert()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			w.Write([]byte("Internal database error"))

		}
		w.WriteHeader(http.StatusCreated)
	}
	// update asset qty
	if r.Method == http.MethodPut {
		var asset database.Asset
		err := json.NewDecoder(r.Body).Decode(&asset)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Malformed asset data"))
		}
		_, er := database.Database.Model(&asset).Set("qty = ?qty").Where("id = ?id").Update()

		if er != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("failed updating asset qty"))
		}
		w.WriteHeader(http.StatusAccepted)

	}
}

func InitServer() {
	log.Println("init server")
	server := http.NewServeMux()
	server.HandleFunc("/", assetsHandler)
	server.HandleFunc("/asset", assetHandler)
	go func() {
		log.Fatal(http.ListenAndServe(port, server)) // need own multiplexer to run in goroutine
	}()
}
