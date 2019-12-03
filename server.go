package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Server type for server
type server struct {
	router *http.ServeMux
}

// NewServer returns new server
func NewServer() *server {
	mux := http.NewServeMux()
	return &server{router: mux}
}

// HandleSignals handles os signals
func (s *server) handleSignals() {
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

func (s *server) defaultHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}

}
