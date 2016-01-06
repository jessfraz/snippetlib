package main

import (
	"flag"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

var (
	dbConn   string
	port     string
	certFile string
	keyFile  string
)

func init() {
	flag.StringVar(&dbConn, "dbconn", "postgres://postgres:@127.0.0.1:5432/db?sslmode=disable", "database connection string")
	flag.StringVar(&port, "p", "3000", "port for server to run on")

	flag.StringVar(&certFile, "cert", "", "path to ssl certificate")
	flag.StringVar(&keyFile, "key", "", "path to ssl key")
	flag.Parse()
}

func main() {
	// create mux router
	r := mux.NewRouter()

	// static files handler
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("/src/static")))
	r.Handle("/static/", staticHandler)

	// template handler
	h := Handler{
		dbConn: dbConn,
	}
	r.HandleFunc("/search", h.sitemapHandler).Methods("GET")
	r.HandleFunc("/search", h.searchHandler).Methods("POST")
	r.HandleFunc("/newsletter_signup", h.newsletterSignupHandler).Methods("POST")
	r.HandleFunc("/{category}/{snippet}", h.snippetHandler).Methods("GET")
	r.HandleFunc("/{category}", h.categoryHandler).Methods("GET")

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
