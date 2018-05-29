package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"log"
	"regexp"
	"errors"
)

type Page struct {
	Title string
	Body  []byte
}

var templatePages = template.Must(template.ParseFiles("edit.html",
	"view.html", ))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

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

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil // The title is the second subexpression.
}

func viewHandler(writer http.ResponseWriter, request *http.Request) {
	title, error := getTitle(writer, request)
	if error != nil {
		return
	}
	page, error := loadPage(title)
	if error != nil {
		http.Redirect(writer, request, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(writer, "view", page)
}

func editHandler(writer http.ResponseWriter, request *http.Request) {
	title, error := getTitle(writer, request)
	if error != nil {
		return
	}
	page, error := loadPage(title)
	if error != nil {
		page = &Page{Title: title}
	}
	renderTemplate(writer, "edit", page)
}

func saveHandler(writer http.ResponseWriter, request *http.Request) {
	title, error := getTitle(writer, request)
	if error != nil {
		return
	}
	body := request.FormValue("body")
	page := &Page{Title: title, Body: []byte(body)}
	page.save()
	http.Redirect(writer, request, "/view/"+title, http.StatusFound)
}

func renderTemplate(writer http.ResponseWriter, templateName string, page *Page) {
	error := templatePages.ExecuteTemplate(writer, templateName+".html", page)
	if error != nil {
		http.Error(writer, error.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
