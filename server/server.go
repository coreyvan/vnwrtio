// Package server implements server library
package server

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Server type for server
type Server struct {
	mux *http.ServeMux
}

// HandleSignals handles os signals
func (s *Server) HandleSignals() {
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
				log.Println("Received shut signal:", s)
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
