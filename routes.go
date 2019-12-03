package main

import (
	"net/http"
)

func (s *server) routes() {
	fs := http.FileServer(http.Dir("static"))
	s.router.Handle("/", fs)
}
