package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	htmlBase   = "./src/templates/"
	cssBase    = "./src/css/"
	jsBase     = "./src/js/"
	assetsBase = "./src/assets/"
	logger     = logrus.New()
)

func init() {

}
func main() {
	certPath := os.Getenv("VNWRT_CERT_PATH")
	privKeyPath := os.Getenv("VNWRT_PRIVKEY_PATH")

	logger.Formatter = new(logrus.TextFormatter)
	logger.Formatter.(*logrus.TextFormatter).FullTimestamp = true

	port := ":443"
	s := NewServer()

	go s.handleSignals()
	s.routes()

	// Create a logger that ignores errors printed by the std library
	ignoreLogger := log.New(ioutil.Discard, "", 0)
	srv := &http.Server{Addr: port, Handler: s.router, ErrorLog: ignoreLogger}

	// Start go routine listening on HTTP port 80 that redirects to HTTPS
	go http.ListenAndServe(":80", http.HandlerFunc(redirectTLS))

	// Start listening HTTPS on port 443
	logger.WithFields(logrus.Fields{
		"package":  "main",
		"function": "main",
		"port":     port,
	}).Info("Started server listening on port", port)
	err := srv.ListenAndServeTLS(certPath, privKeyPath)
	if err != nil {
		logger.WithFields(logrus.Fields{
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
	s.router.Handle("/{page:contact|sharks}", s.pageHandler()).Methods("GET")
	s.router.Handle("/{res:css|js}/{file}", s.resourceHandler()).Methods("GET")
	s.router.Handle("/favicon.ico", s.faviconHandler()).Methods("GET")
	s.router.Handle("/", s.badSchemeHandler()).Schemes("http")
	s.router.Handle("/", s.homeHandler()).Methods("GET")
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
			processSignal(s, exitChan)
		}
	}()

	code := <-exitChan
	os.Exit(code)
}

func processSignal(s os.Signal, exit chan int) {
	switch s {
	case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
		logger.WithFields(logrus.Fields{
			"package":  "main",
			"function": "processSignal",
			"signal":   s,
		}).Info("Received shutdown signal")
		exit <- 0
	default:
		logger.WithFields(logrus.Fields{
			"package":  "main",
			"function": "processSignals",
			"signal":   s,
		}).Error("Received unknown signal")
		exit <- 1
	}
}

func (s *Server) badSchemeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	}
}

// Page data model for HTML page
type Page struct {
	Title string
}

func (s *Server) homeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		logger.WithFields(logrus.Fields{
			"package":  "main",
			"function": "homeHandler",
			"method":   r.Method,
			"source":   r.RemoteAddr,
		}).Info("Request to ", r.URL.Path)
		p := Page{Title: "Corey Van Woert"}

		// Parse HTML template
		tmpl, err := template.ParseFiles(htmlBase + "index.html")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			logger.WithFields(logrus.Fields{
				"package":  "main",
				"function": "homeHandler",
				"template": htmlBase + "index.html",
				"source":   r.RemoteAddr,
			}).Errorf("Could not execute template: %v", err)
		}

		// Write to client
		w.Header().Add("Access-Control-Allow-Origin", "*")
		renderPageTemplate(w, tmpl, p)

	}
}

func (s *Server) resourceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.WithFields(logrus.Fields{
			"package":  "main",
			"function": "resourceHandler",
			"method":   r.Method,
			"source":   r.RemoteAddr,
		}).Info("Request to ", r.URL.Path)

		// Get variables from path, res will only be an item in the patterns matched in the route
		vars := mux.Vars(r)
		t := vars["res"]
		f := vars["file"]
		if containsDotDot(f) {
			// Reject requests with ..'s in the path to avoid directory traversal attacks
			// technically http.ServeFile does this by default but don't take a dependency on it
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

			logger.WithFields(logrus.Fields{
				"package":  "main",
				"function": "resourceHandler",
				"method":   r.Method,
				"source":   r.RemoteAddr,
				"URL":      r.URL.Path,
			}).Error("URL path contained ..")
		}

		// Serve file based on captured path
		http.ServeFile(w, r, path.Join("src", t, f))
	}
}

func containsDotDot(v string) bool {
	if !strings.Contains(v, "..") {
		return false
	}
	for _, ent := range strings.FieldsFunc(v, isSlashRune) {
		if ent == ".." {
			return true
		}
	}
	return false
}

func isSlashRune(r rune) bool { return r == '/' || r == '\\' }

func (s *Server) pageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.WithFields(logrus.Fields{
			"package":  "main",
			"function": "pageHandler",
			"method":   r.Method,
			"source":   r.RemoteAddr,
		}).Info("Request to ", r.URL.Path)

		vars := mux.Vars(r)
		page := vars["page"]
		path := htmlBase + page + ".html"
		title := fmt.Sprintf("%s - Corey Van Woert", strings.Title(page))

		p := Page{Title: title}

		tmpl, err := template.ParseFiles(path)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			logger.WithFields(logrus.Fields{
				"package":  "main",
				"function": "pageHandler",
				"template": path,
				"source":   r.RemoteAddr,
			}).Errorf("Could not execute template: %v", err)
			return
		}

		w.Header().Add("Access-Control-Allow-Origin", "*")
		renderPageTemplate(w, tmpl, p)
	}
}

func (s *Server) faviconHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := assetsBase + "favicon.ico"
		http.ServeFile(w, r, path)
	}
}
func redirectTLS(w http.ResponseWriter, r *http.Request) {
	logger.WithFields(logrus.Fields{
		"package":  "main",
		"function": "redirectTLS",
		"method":   r.Method,
		"source":   r.RemoteAddr,
	}).Info("HTTP request redirected to HTTPS")
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
}

func renderPageTemplate(w http.ResponseWriter, tmpl *template.Template, p Page) {
	err := tmpl.Execute(w, p)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		logger.WithFields(logrus.Fields{
			"package":  "main",
			"function": "renderPageTemplate",
			"template": tmpl,
			"err":      err,
		}).Error("Error while executing template")
	}
}
