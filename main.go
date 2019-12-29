package main

import (
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	htmlBase   = "./src/templates/"
	cssBase    = "./src/css/"
	jsBase     = "./src/js/"
	assetsBase = "./src/assets/"
	log        = logrus.New()
)

func init() {

}
func main() {
	certPath := os.Getenv("VNWRT_CERT_PATH")
	privKeyPath := os.Getenv("VNWRT_PRIVKEY_PATH")

	log.Formatter = new(logrus.TextFormatter)
	log.Formatter.(*logrus.TextFormatter).FullTimestamp = true

	port := ":443"
	s := NewServer()

	go s.handleSignals()
	s.routes()

	go http.ListenAndServe(":80", http.HandlerFunc(redirectTLS))
	log.WithFields(logrus.Fields{
		"package":  "main",
		"function": "main",
		"port":     port,
	}).Info("Started server")
	err := http.ListenAndServeTLS(port, certPath, privKeyPath, s.router)
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
	router *mux.Router
}

// NewServer returns new server
func NewServer() *Server {
	mux := mux.NewRouter()
	return &Server{router: mux}
}

// Define routes
func (s *Server) routes() {
	s.router.Handle("/contact", s.contactHandler()).Methods("GET")
	s.router.Handle("/", s.homeHandler()).Methods("GET")
	s.router.Handle("/{res:css|js}/{file}", s.resourceHandler()).Methods("GET")
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
}

func (s *Server) homeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(logrus.Fields{
			"package":  "main",
			"function": "homeHandler",
			"method":   r.Method,
			"source":   r.RemoteAddr,
		}).Info("Request to ", r.URL.Path)
		p := Page{Title: "Corey Van Woert"}
		tmpl, err := template.ParseFiles(htmlBase + "index.html")
		if err != nil {
			log.WithFields(logrus.Fields{
				"package":  "main",
				"function": "homeHandler",
				"template": htmlBase + "index.html",
				"source":   r.RemoteAddr,
			}).Errorf("Could not execute template: %v", err)
		}
		w.Header().Add("Access-Control-Allow-Origin", "*")
		renderPageTemplate(w, tmpl, p)

	}
}

func (s *Server) resourceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(logrus.Fields{
			"package":  "main",
			"function": "resourceHandler",
			"method":   r.Method,
			"source":   r.RemoteAddr,
		}).Info("Request to ", r.URL.Path)
		vars := mux.Vars(r)
		t := vars["res"]
		f := vars["file"]
		http.ServeFile(w, r, path.Join("src", t, f))
	}
}

func (s *Server) contactHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(logrus.Fields{
			"package":  "main",
			"function": "contactHandler",
			"method":   r.Method,
		}).Info("Request to ", r.URL.Path)
		p := Page{Title: "Contact - Corey Van Woert"}
		tmpl, err := template.ParseFiles(htmlBase + "contact.html")
		if err != nil {
			log.WithFields(logrus.Fields{
				"package":  "main",
				"function": "contactHandler",
				"template": htmlBase + "contact.html",
				"source":   r.RemoteAddr,
			}).Errorf("Could not execute template: %v", err)
		}
		w.Header().Add("Access-Control-Allow-Origin", "*")
		renderPageTemplate(w, tmpl, p)
	}
}

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
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
