package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/jessfraz/snippetlib/mailchimp"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
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

func (h *Handler) searchHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("[page] %s", r.URL)

	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseForm(); err != nil {
		writeError(w, r, fmt.Sprintf("parsing form failed: %v", err))
		return
	}

	data, err := search(h.dbConn, r.Form.Get("category"), r.Form.Get("q"))
	if err != nil {
		writeError(w, r, err.Error())
		return
	}
	str, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		logrus.Warnf("marshaling search data failed: %v", err)
	}

	fmt.Fprint(w, string(str))
}

func (h *Handler) newsletterSignupHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("[page] %s", r.URL)

	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseForm(); err != nil {
		writeError(w, r, fmt.Sprintf("parsing form failed: %v", err))
		return
	}

	email := r.Form.Get("email")
	if err := mailchimp.Subscribe(mailchimpAPIKey, mailchimpListID, email); err != nil {
		writeError(w, r, fmt.Sprintf("subscribing email (%s) to list (%s) failed: %v", email, mailchimpListID, err))
		return
	}

	logrus.Infof("subscribed %s to newsletter", email)
	fmt.Fprint(w, JSONResponse{
		"response": "success",
	})
}

func (h *Handler) categoryHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("[page] %s", r.URL)

	v := mux.Vars(r)
	category, ok := v["category"]
	if !ok {
		writeError(w, r, fmt.Sprintf("getting category parameter from vars failed: %v", v))
		return
	}

	h.renderTemplate(w, r, category, "")
}

func (h *Handler) snippetHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("[page] %s", r.URL)

	v := mux.Vars(r)
	category, ok := v["category"]
	if !ok {
		writeError(w, r, fmt.Sprintf("getting category parameter from vars failed: %v", v))
		return
	}
	slug, ok := v["snippet"]
	if !ok {
		writeError(w, r, fmt.Sprintf("getting snippet parameter from vars failed: %v", v))
		return
	}

	h.renderTemplate(w, r, category, slug)
}

func (h *Handler) indexHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("[page] %s", r.URL)

	h.renderTemplate(w, r, "facebook", "")
}

func (h *Handler) renderTemplate(w http.ResponseWriter, r *http.Request, category, slug string) {
	page, err := query(h.dbConn, category, slug)
	if err != nil {
		writeError(w, r, err.Error())
		return
	}
	page.URL = r.URL.String()

	if slug == "" && len(page.Snippets) == 0 {
		writeError(w, r, fmt.Sprintf("no snippets found for category (%s)", category))
		return
	}

	// render the template
	lp := path.Join(templateDir, "layout.html")

	// set up custom functions
	funcMap := template.FuncMap{
		"stripSlashes": func(s string) string {
			return strings.Replace(s, "/", "", -1)
		},
		"stripTags": func(s string) string {
			reHTML := regexp.MustCompile(`<.+>`)
			return reHTML.ReplaceAllString(s, "")
		},
		"toTitle": func(s string) string {
			return strings.ToTitle(s)
		},
		"toUppercase": func(s string) string {
			return strings.ToUpper(s)
		},
	}

	// parse & execute the template
	tmpl := template.Must(template.New("").Funcs(funcMap).ParseFiles(lp))
	if err := tmpl.ExecuteTemplate(w, "layout", page); err != nil {
		writeError(w, r, fmt.Sprintf("execute template failed: %v", err))
		return
	}
}

// writeError logs the error and redirects the user to /.
func writeError(w http.ResponseWriter, r *http.Request, msg string) {
	logrus.Errorf("handler error: %s", msg)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
