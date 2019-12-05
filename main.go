package main

import (
	"html/template"
	"log"
	"net/http"
	"regexp"
)

var (
	templates = template.Must(template.ParseFiles("templates/view.html", "templates/edit.html", "templates/index.html"))
	validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
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
