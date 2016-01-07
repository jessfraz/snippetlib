package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

const (
	filesPrefix = "/src"
)

var (
	dbConn   string
	port     string
	certFile string
	keyFile  string
	debug    bool

	templateDir = path.Join(filesPrefix, "templates")
)

func init() {
	flag.StringVar(&dbConn, "dbconn", "postgres://postgres:@127.0.0.1:5432/db?sslmode=disable", "database connection string")
	flag.StringVar(&port, "p", "3000", "port for server to run on")

	flag.StringVar(&certFile, "cert", "", "path to ssl certificate")
	flag.StringVar(&keyFile, "key", "", "path to ssl key")

	flag.BoolVar(&debug, "d", false, "run in debug mode")
	flag.Parse()

	// set log level
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func main() {
	// get the sitemap
	sitemap, err := getSitemap(dbConn)
	if err != nil {
		logrus.Fatal(err)
	}

	// create mux router
	r := mux.NewRouter()
	r.StrictSlash(true)

	// static files handler
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir(path.Join(filesPrefix, "static"))))
	r.PathPrefix("/static/").Handler(staticHandler)

	// template handler
	h := Handler{
		dbConn:  dbConn,
		sitemap: sitemap,
	}
	r.HandleFunc("/sitemap.xml", h.sitemapHandler).Methods("GET")
	r.HandleFunc("/search", h.searchHandler).Methods("POST")
	r.HandleFunc("/newsletter_signup", h.newsletterSignupHandler).Methods("POST")
	r.HandleFunc("/{category}", h.categoryHandler).Methods("GET")
	r.HandleFunc("/{category}/{snippet}", h.snippetHandler).Methods("GET")
	r.HandleFunc("/", h.indexHandler).Methods("GET")

	// set up the server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	logrus.Infof("Starting server on port %q", port)
	if certFile != "" && keyFile != "" {
		logrus.Fatal(server.ListenAndServeTLS(certFile, keyFile))
	} else {
		logrus.Fatal(server.ListenAndServe())
	}
}

func getSitemap(dbConn string) ([]byte, error) {
	urls, err := sitemapQuery(dbConn)
	if err != nil {
		return nil, err
	}

	// render the template
	sm := path.Join(templateDir, "sitemap.xml")
	tmpl := template.Must(template.New("").ParseFiles(sm))

	// parse & execute the template
	var b bytes.Buffer
	if err := tmpl.ExecuteTemplate(bufio.NewWriter(&b), "sitemap", urls); err != nil {
		return nil, fmt.Errorf("execute sitemap template failed: %v", err)
	}

	return b.Bytes(), err
}
