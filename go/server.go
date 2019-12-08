package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Server type for server
type Server struct {
	router *http.ServeMux
}

// NewServer returns new server
func NewServer() *Server {
	mux := http.NewServeMux()
	return &Server{router: mux}
}

// HandleSignals handles os signals
func (s *Server) handleSignals() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	exitChan := make(chan int)

	go func() {
		for {
			s := <-signalChan
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				log.Println("Received shutdown signal:", s)
				exitChan <- 0
			default:
				log.Println("Received unknown signal:", s)
				exitChan <- 1
			}
		}
	}()

	code := <-exitChan
	os.Exit(code)
}

func (s *Server) defaultHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}

}

type Page struct {
	Title string
	Cards []Card
}

type Card struct {
	ID string
	Text string
	Link Link
}

type Link struct {
	Text string
	Link string
}

func (s *Server) homeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cards := []Card{
			Card{"1", "card 1 text", Link{"text", "link"}},
			Card{"2", "card 2 text", Link{"text", "link"}},
			Card{"3", "card 3 text", Link{"text", "link"}},
		}
		p := Page{Title: "Corey Van Woert", Cards: cards}
		renderPageTemplate(w, "index", p)
	}
}

func renderPageTemplate(w http.ResponseWriter, tmpl string, p Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
