package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"text/template"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view|createNewPage|doCreate)/([a-zA-Z0-9]+)$")

type Page struct {
	Title string
	Route string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"

	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

func getTitle(w http.ResponseWriter, r *http.Request) string {
	paths := validPath.FindStringSubmatch(r.URL.Path)
	if paths == nil {
		http.NotFound(w, r)
	}
	return paths[2]
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := getTitle(w, r)
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/createNewPage/"+title, http.StatusFound)
		return
	}

	p.Route = "view"
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := getTitle(w, r)
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/createNewPage/"+title, http.StatusFound)
		return
	}

	p.Route = "edit"
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := getTitle(w, r)
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func createPageHandler(w http.ResponseWriter, r *http.Request) {
	title := getTitle(w, r)
	p := &Page{Title: title, Route: "createPage"}

	renderTemplate(w, "edit", p)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	fmt.Println("Program is running...")

	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/createNewPage/", createPageHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
