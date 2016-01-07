package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	defaultDescription = "A collection of coding language and api snippets. Code snippets for html, css, javascript, jquery, php, sql, and facebook, instagram, twitter, google analytics, google maps, google charts, foursquare, wordpress, tumblr, last.fm and pinterest api."
)

// Page defines the struct that gets passed to the layout template.
type Page struct {
	Category    Category
	Count       int
	Description string
	Snippet     Snippet
	Snippets    []*Snippet
	URL         string
}

// Category is how snippets get organized into groups.
type Category struct {
	Greeting string
	Name     string
	Slug     string
}

// Snippet contains the code and description for a code snippet.
type Snippet struct {
	Category string
	Name     string
	Slug     string
	Snippet  string
}

func query(dbConn, categoryPassed, snippet string) (p Page, err error) {
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		return p, fmt.Errorf("opening database at %s failed: %v", dbConn, err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT code_snippets.category AS category, code_snippets.name AS name, code_snippets.slug AS slug, code_categories.name AS category_formal, code_categories.header AS greeting FROM code_snippets LEFT JOIN code_categories ON code_snippets.category=code_categories.category WHERE code_snippets.category='" + categoryPassed + "' ORDER BY name ASC")
	if err != nil {
		return p, fmt.Errorf("querying for category (%s) failed: %v", categoryPassed, err)
	}
	defer rows.Close()

	for rows.Next() {
		var category, categoryFormal, greeting, name, slug string
		if err := rows.Scan(&category, &categoryFormal, &greeting, &name, &slug); err != nil {
			return p, fmt.Errorf("scanning rows for category (%s) and fields failed: %v", categoryPassed, err)
		}
		fmt.Printf("got %s, %s, %s\n", category, slug, name)
		p.Snippets = append(p.Snippets, &Snippet{
			Category: category,
			Name:     name,
			Slug:     slug,
		})
		p.Category = Category{
			Greeting: greeting,
			Name:     categoryFormal,
			Slug:     categoryPassed,
		}
	}
	if err := rows.Err(); err != nil {
		return p, fmt.Errorf("scanning rows for category (%s) overall failed: %v", categoryPassed, err)
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
				Category: category,
				Name:     name,
				Slug:     slug,
				Snippet:  snippet,
			}
			p.Description = description
		}
	}

	return p, nil
}

func search(dbConn, categoryPassed, q string) ([]map[string]interface{}, error) {
	query := "SELECT category, name, slug FROM code_snippets WHERE category='" + categoryPassed + "'"
	if q != "" {
		query += " AND name LIKE '%" + q + "%' ORDER BY name ASC"
	} else {
		query += " ORDER BY name ASC"
	}

	return jsonQuery(dbConn, query)
}

func sitemapQuery(dbConn string) ([]map[string]interface{}, error) {
	q := "SELECT code_snippets.category AS category, code_snippets.slug AS slug FROM code_snippets"

	return jsonQuery(dbConn, q)
}

func jsonQuery(dbConn, q string) (data []map[string]interface{}, err error) {
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		return data, fmt.Errorf("opening database at %s failed: %v", dbConn, err)
	}
	defer db.Close()

	rows, err := db.Query(q)
	if err != nil {
		return data, fmt.Errorf("searching for query (%s) failed: %v", q, err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return data, fmt.Errorf("getting columns for query (%s) failed: %v", q, err)
	}
	count := len(columns)

	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		data = append(data, entry)
	}
	return data, nil
}
