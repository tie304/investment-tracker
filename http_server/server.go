package http_server

import (
	"log"
	"net/http"
)

const port = ":8000"

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("using handler")
}

func InitServer() {
	log.Println("init server")
	server := http.NewServeMux()
	http.HandleFunc("/", handler)
	go func() {
		http.ListenAndServe(port, server) // need own multiplexer to run in goroutine
	}()
}
