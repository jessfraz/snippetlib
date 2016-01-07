package main

import (
	"fmt"
	"os"
	"path"
	"text/template"
)

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
