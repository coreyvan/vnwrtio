package main

import (
	"html/template"
	"log"
	"net/http"
)

var (
	htmlBase  = "../src/"
	cssBase   = "../src/css/"
	jsBase    = "../src/js/"
	templates = template.Must(template.ParseFiles(htmlBase + "index.html"))
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
