package main

import (
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

var (
	// htmlBase  = "../src/"
	// cssBase   = "../src/css/"
	// jsBase    = "../src/js/"
	htmlBase   = "./src/templates/"
	cssBase    = "./src/css/"
	jsBase     = "./src/js/"
	assetsBase = "./src/assets/"
	// templates  = template.Must(template.ParseFiles(htmlBase + "index.html"))
	log = logrus.New()
)

func init() {

}
func main() {

	log.Formatter = new(logrus.TextFormatter)
	log.Formatter.(*logrus.TextFormatter).FullTimestamp = true

	port := ":80"
	s := NewServer()

	go s.handleSignals()
	s.routes()

	log.WithFields(logrus.Fields{
		"package":  "main",
		"function": "main",
		"port":     port,
	}).Info("Started server")
	err := http.ListenAndServe(port, s.router)
	if err != nil {
		log.WithFields(logrus.Fields{
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
	s.router.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(assetsBase))))
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
				log.WithFields(logrus.Fields{
					"package":  "main",
					"function": "handleSignals",
					"signal":   s,
				}).Info("Received shutdown signal")
				exitChan <- 0
			default:
				log.WithFields(logrus.Fields{
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
	Title string
	List  []string
}

// Link data model for HTML link
type Link struct {
	Text string
	Link string
}

func (s *Server) homeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(logrus.Fields{
			"package":  "main",
			"function": "homeHandler",
			"method":   r.Method,
		}).Info("Request to ", r.URL.Path)
		cards := []Card{
			Card{"Code", []string{"Go"}},
			Card{"Deploy", []string{"Docker", "Terraform"}},
		}
		p := Page{Title: "Corey Van Woert", Cards: cards}
		tmpl, err := template.ParseFiles(htmlBase + "index.html")
		if err != nil {
			log.WithFields(logrus.Fields{
				"package":  "main",
				"function": "homeHandler",
				"template": htmlBase + "index.html",
			}).Errorf("Could not execute template: %v", err)
		}
		w.Header().Add("Access-Control-Allow-Origin", "*")
		renderPageTemplate(w, tmpl, p)

	}
}

func renderPageTemplate(w http.ResponseWriter, tmpl *template.Template, p Page) {
	err := tmpl.Execute(w, p)
	if err != nil {
		log.WithFields(logrus.Fields{
			"package":  "main",
			"function": "renderPageTemplate",
			"template": tmpl,
			"err":      err,
		}).Error("Error while executing template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
