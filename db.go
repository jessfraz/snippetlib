package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	defaultDescription = "A collection of coding language and api snippets. Code snippets for html, css, javascript, jquery, php, sql, and facebook, instagram, twitter, google analytics, google maps, google charts, foursquare, wordpress, tumblr, last.fm and pinterest api."
)

type Page struct {
	Category    Category
	Count       int
	Description string
	Snippet     Snippet
	Snippets    []Snippet
	URL         string
}

type Category struct {
	greeting string
	name     string
	slug     string
}

type Snippet struct {
	category string
	name     string
	slug     string
	snippet  string
}

func query(dbConn, categoryPassed, snippet string) (p Page, err error) {
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		return p, fmt.Errorf("opening database at %s failed: %v", dbConn, err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT code_snippets.category AS cat,
	code_snippets.name AS name,
	code_snippets.slug AS slug,
	code_categories.name AS categoryFormal,
	code_categories.header AS greeting
	FROM code_snippets LEFT JOIN code_categories
	ON code_snippets.category=code_categories.category
	WHERE code_snippets.category=$1 ORDER BY name ASC`, categoryPassed)
	if err != nil {
		return p, fmt.Errorf("querying for category (%s) failed: %v", categoryPassed, err)
	}
	defer rows.Close()

	var categoryFormal, greeting string
	for rows.Next() {
		var category, name, slug string
		if err := rows.Scan(&category, &categoryFormal, &greeting, &name, &slug); err != nil {
			return p, fmt.Errorf("scanning rows for category (%s) and fields failed: %v", categoryPassed, err)
		}
		p.Snippets = append(p.Snippets, Snippet{
			category: category,
			name:     name,
			slug:     slug,
		})
	}
	if err := rows.Err(); err != nil {
		return p, fmt.Errorf("scanning rows for category (%s) overall failed: %v", categoryPassed, err)
	}
	p.Category = Category{
		greeting: greeting,
		name:     categoryFormal,
		slug:     categoryPassed,
	}
	p.Description = defaultDescription

	if snippet != "" {
		var category, description, name, slug, snippet string
		err = db.QueryRow(`SELECT * FROM code_snippets WHERE slug=$1 AND category=$2 LIMIT 1`, snippet, categoryPassed).Scan(&category, &description, &name, &slug, &snippet)
		switch {
		case err == sql.ErrNoRows:
			return p, fmt.Errorf("querying for category (%s) and snippet (%s) returned no rows", categoryPassed, snippet)
		case err != nil:
			return p, fmt.Errorf("querying for category (%s) and snippet (%s) failed: %v", categoryPassed, snippet, err)
		default:
			p.Snippet = Snippet{
				category: category,
				name:     name,
				slug:     slug,
				snippet:  snippet,
			}
			p.Description = description
		}
	}

	return p, nil
}
