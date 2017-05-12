package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var templates map[string]*template.Template

// Load templates on program initialisation
func init() {
	templates = make(map[string]*template.Template)
	templates["index.html"] = template.Must(template.ParseFiles("templates/index.html", "templates/nav.html", "templates/base.html"))
	templates["upload.html"] = template.Must(template.ParseFiles("templates/upload.html", "templates/nav.html", "templates/base.html"))
	templates["login.html"] = template.Must(template.ParseFiles("templates/login.html", "templates/nav.html", "templates/base.html"))
}

// renderTemplate is a wrapper around template.ExecuteTemplate.
func renderTemplate(w http.ResponseWriter, name string, data map[string]interface{}) error {
	// Ensure the template exists in the map.
	tmpl, ok := templates[name]
	if !ok {
		return fmt.Errorf("The template %s does not exist.", name)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return tmpl.ExecuteTemplate(w, "base", data)
}
