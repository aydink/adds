package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

var templates map[string]*template.Template

var categoryMap map[int]string
var regionMap map[int]string

// Load templates on program initialisation
func init() {

	funcMap := template.FuncMap{
		"category": GetCategory,
		"region":   GetRegion,
		"newline":  NewLineToBr,
		"path":     GetPath,
	}

	templates = make(map[string]*template.Template)
	templates["index.html"] = template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/index.html", "templates/nav.html", "templates/base.html"))
	templates["upload.html"] = template.Must(template.ParseFiles("templates/upload.html", "templates/nav.html", "templates/base.html"))
	templates["login.html"] = template.Must(template.ParseFiles("templates/login.html", "templates/nav.html", "templates/base.html"))
	templates["new_edit_add.html"] = template.Must(template.ParseFiles("templates/new_edit_add.html", "templates/nav.html", "templates/base.html"))
	templates["adds_myadds.html"] = template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/adds_myadds.html", "templates/nav.html", "templates/base.html"))
	//templates["upload_photo.html"] = template.Must(template.ParseFiles("templates/upload_photo.html", "templates/nav.html", "templates/base.html")).Funcs(funcMap)
	templates["upload_photo.html"] = template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/upload_photo.html", "templates/nav.html", "templates/base.html"))
	templates["adds_view.html"] = template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/adds_view.html", "templates/nav.html", "templates/base.html"))

	// init category and region maps
	categoryMap = map[int]string{
		1: "Mobilya",
		2: "Oto",
		3: "Beyaz eşya",
		4: "Bakıcı/Temizlikçi",
	}

	regionMap = map[int]string{
		1: "ORAN",
		2: "ÇANKAYA",
		3: "ÇİĞİLTEPE",
		4: "ÇAĞLAYAN",
		5: "ERLER MAHALLESİ",
		6: "ÇAĞLAYAN",
	}

	//templates["upload_photo.html"].Funcs(funcMap)
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

func GetCategory(id int) string {
	return categoryMap[id]
}

func GetRegion(id int) string {
	return regionMap[id]
}

func NewLineToBr(text string) template.HTML {
	return template.HTML(strings.Replace(text, "\n", "<br>", -1))
}

func GetPath(hash string) string {
	if len(hash) > 2 {
		return hash[0:2]
	}

	return ""
}
