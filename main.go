package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"text/template"

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
	_, err := getSitemap(dbConn)
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
		dbConn: dbConn,
	}
	r.PathPrefix("/sitemap.xml").Handler(http.StripPrefix("/", http.FileServer(http.Dir(filesPrefix))))
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

// getSitemap takes a database connection and generates a sitemap from it.
func getSitemap(dbConn string) (string, error) {
	urls, err := sitemapQuery(dbConn)
	if err != nil {
		return "", err
	}

	// render the template
	sm := path.Join(templateDir, "sitemap.xml")
	tmpl := template.Must(template.New("").ParseFiles(sm))

	// open the file we will write to
	sitemap := path.Join(filesPrefix, "sitemap.xml")
	f, err := os.OpenFile(sitemap, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return "", fmt.Errorf("opening %s for sitemap failed: %v", sitemap, err)
	}
	defer f.Close()

	// parse & execute the template
	if err := tmpl.ExecuteTemplate(f, "sitemap", urls); err != nil {
		return "", fmt.Errorf("execute sitemap template failed: %v", err)
	}

	return sitemap, err
}
