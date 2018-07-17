package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/genuinetools/pkg/cli"
	"github.com/gorilla/mux"
	"github.com/jessfraz/snippetlib/version"
	"github.com/sirupsen/logrus"
)

const (
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

	templateDir = path.Join(filesPrefix, "templates")
)

func main() {
	// Create a new cli program.
	p := cli.NewProgram()
	p.Name = "snippetlib"
	p.Description = "Server to host code snippets"

	// Set the GitCommit and Version.
	p.GitCommit = version.GITCOMMIT
	p.Version = version.VERSION

	// Setup the global flags.
	p.FlagSet = flag.NewFlagSet("global", flag.ExitOnError)
	p.FlagSet.StringVar(&dbConn, "dbconn", "postgres://postgres:@127.0.0.1:5432/db?sslmode=disable", "database connection string")
	p.FlagSet.StringVar(&mailchimpAPIKey, "mailchimp-apikey", "", "Mailchimp APIKey for subscribing to newsletters")
	p.FlagSet.StringVar(&mailchimpListID, "mailchimp-listid", "", "Mailchimp List ID for newsletter to subscribe emails to")

	p.FlagSet.StringVar(&port, "p", "3000", "port for server to run on")
	p.FlagSet.StringVar(&certFile, "cert", "", "path to ssl certificate")
	p.FlagSet.StringVar(&keyFile, "key", "", "path to ssl key")

	p.FlagSet.BoolVar(&debug, "d", false, "enable debug logging")

	// Set the before function.
	p.Before = func(ctx context.Context) error {
		// Set the log level.
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		return nil
	}

	// Set the main program action.
	p.Action = func(ctx context.Context, args []string) error {
		// On ^C, or SIGTERM handle exit.
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		signal.Notify(c, syscall.SIGTERM)
		go func() {
			for sig := range c {
				logrus.Infof("Received %s, exiting.", sig.String())
				os.Exit(0)
			}
		}()

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

		return nil
	}

	// Run our program.
	p.Run()
}
