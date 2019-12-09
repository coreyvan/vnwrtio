package main

import (
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

var (
	htmlBase  = "../src/"
	cssBase   = "../src/css/"
	jsBase    = "../src/js/"
	templates = template.Must(template.ParseFiles(htmlBase + "index.html"))
)

func init() {

}
func main() {
	port := ":80"
	s := NewServer()

	go s.handleSignals()
	s.routes()

	log.WithFields(log.Fields{
		"package":  "main",
		"function": "main",
		"port":     port,
	}).Info("Started server")
	err := http.ListenAndServe(port, s.router)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "main",
			"function": "main",
			"err":      err,
		}).Error("Encountered error while listening")
	}
}

// Server type for server
type Server struct {
	router *http.ServeMux
}

// NewServer returns new server
func NewServer() *Server {
	mux := http.NewServeMux()
	return &Server{router: mux}
}

func (s *Server) routes() {
	// fs := http.FileServer(http.Dir("static"))
	s.router.Handle("/", s.homeHandler())
	s.router.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(cssBase))))
	s.router.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(jsBase))))
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
				log.WithFields(log.Fields{
					"package":  "main",
					"function": "handleSignals",
					"signal":   s,
				}).Info("Received shutdown signal")
				exitChan <- 0
			default:
				log.WithFields(log.Fields{
					"package":  "main",
					"function": "handleSignals",
					"signal":   s,
				}).Error("Received unknown signal")
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

// Page data model for HTML page
type Page struct {
	Title string
	Cards []Card
}

// Card data model for HTML card
type Card struct {
	ID   string
	Text string
	Link Link
}

// Link data model for HTML link
type Link struct {
	Text string
	Link string
}

func (s *Server) homeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"package":  "main",
			"function": "homeHandler",
			"method":   r.Method,
		}).Info("Request to", r.URL.Path)
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
		log.WithFields(log.Fields{
			"package":  "main",
			"function": "renderPageTemplate",
			"template": tmpl,
			"err":      err,
		}).Error("Error while executing template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
