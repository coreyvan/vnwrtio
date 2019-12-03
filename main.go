package main

import (
	"log"
	"net/http"
)

func main() {
	port := ":80"
	s := NewServer()

	go s.handleSignals()
	s.routes()

	log.Println("Started server listening on port", port)
	err := http.ListenAndServe(port, s.router)
	if err != nil {
		log.Printf("error while listening: %v", err)
	}
}
