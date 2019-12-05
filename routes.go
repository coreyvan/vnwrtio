package main

import "net/http"

func (s *Server) routes() {
	// fs := http.FileServer(http.Dir("static"))
	s.router.Handle("/", s.homeHandler())
	s.router.Handle("/view/", s.viewHandler())
	s.router.Handle("/edit/", s.editHandler())
	s.router.Handle("/save/", s.saveHandler())
	s.router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}
