package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// JSONResponse is a map[string]string
// response from the web server.
type JSONResponse map[string]string

// String returns the string representation of the
// JSONResponse object.
func (j JSONResponse) String() string {
	str, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{
  "error": "%v"
}`, err)
	}

	return string(str)
}

// Handler is the object which contains data to pass to the http handler functions.
type Handler struct {
	dbConn string
}

func (h *Handler) sitemapHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, JSONResponse{
		"page": "sitemap.xml",
	})
}

func (h *Handler) searchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, JSONResponse{
		"page": "search",
	})
}

func (h *Handler) newsletterSignupHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, JSONResponse{
		"page": "newsletter_signup",
	})
}

func (h *Handler) categoryHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	category, ok := v["category"]
	if !ok {
		writeError(w, fmt.Sprintf("getting category parameter from vars failed: %#v", v))
	}

	page, err := query(h.dbConn, category, "")
	if err != nil {
		writeError(w, err.Error())
	}
	page.URL = r.URL.String()
	logrus.Infof("result: %#v", page)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, JSONResponse{
		"page": "category",
	})
}

func (h *Handler) snippetHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	category, ok := v["category"]
	if !ok {
		writeError(w, fmt.Sprintf("getting category parameter from vars failed: %#v", v))
	}
	slug, ok := v["snippet"]
	if !ok {
		writeError(w, fmt.Sprintf("getting snippetparameter from vars failed: %#v", v))
	}

	page, err := query(h.dbConn, category, slug)
	if err != nil {
		writeError(w, err.Error())
	}
	page.URL = r.URL.String()
	logrus.Infof("result: %#v", page)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, JSONResponse{
		"page": "category/snippet",
	})
}

// writeError sends an error back to the requester
// and also logs the error.
func writeError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, JSONResponse{
		"error": msg,
	})
	logrus.Printf("writing error: %s", msg)
	return
}
