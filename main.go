package main

import (
	"log"
	"net/http"

	"github.com/coreyvan/vnwrtio/server"
)

func main() {
	s := server.Server{}

	go s.HandleSignals()

	log.Println("Started server listening on port 8000")
	http.ListenAndServe(":8000", s.Mux)
}
