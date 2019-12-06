package main

import "net/http"

func (s *Server) routes() {
	// fs := http.FileServer(http.Dir("static"))
	s.router.Handle("/", s.homeHandler())
	s.router.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(cssBase))))
	s.router.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(jsBase))))
}
