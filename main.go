package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
	"github.com/jessfraz/snippetlib/version"
	"github.com/sirupsen/logrus"
)

const (
	// BANNER is what is printed for help/info output.
	BANNER = `           _                  _   _ _ _
 ___ _ __ (_)_ __  _ __   ___| |_| (_) |__
/ __| '_ \| | '_ \| '_ \ / _ \ __| | | '_ \
\__ \ | | | | |_) | |_) |  __/ |_| | | |_) |
|___/_| |_|_| .__/| .__/ \___|\__|_|_|_.__/
            |_|   |_|

 Server to host code snippets.
 Version: %s
 Build: %s

`

	filesPrefix = "/src"
)

var (
	dbConn          string
	mailchimpAPIKey string
	mailchimpListID string

	port     string
	certFile string
	keyFile  string
	debug    bool
	vrsn     bool

	templateDir = path.Join(filesPrefix, "templates")
)

func init() {
	flag.StringVar(&dbConn, "dbconn", "postgres://postgres:@127.0.0.1:5432/db?sslmode=disable", "database connection string")
	flag.StringVar(&mailchimpAPIKey, "mailchimp-apikey", "", "Mailchimp APIKey for subscribing to newsletters")
	flag.StringVar(&mailchimpListID, "mailchimp-listid", "", "Mailchimp List ID for newsletter to subscribe emails to")

	flag.StringVar(&port, "p", "3000", "port for server to run on")
	flag.StringVar(&certFile, "cert", "", "path to ssl certificate")
	flag.StringVar(&keyFile, "key", "", "path to ssl key")

	flag.BoolVar(&vrsn, "version", false, "print version and exit")
	flag.BoolVar(&vrsn, "v", false, "print version and exit (shorthand)")
	flag.BoolVar(&debug, "d", false, "run in debug mode")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER, version.VERSION, version.GITCOMMIT))
		flag.PrintDefaults()
	}

	flag.Parse()

	if vrsn {
		fmt.Printf("snippetlib version %s, build %s", version.VERSION, version.GITCOMMIT)
		os.Exit(0)
	}

	// set log level
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func main() {
	// get the sitemap
	if _, err := getSitemap(dbConn); err != nil {
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
