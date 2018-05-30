package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"log"
	"regexp"
	"fmt"
)

type Page struct {
	Title string
	Body  []byte
}

var templatePages = template.Must(template.ParseFiles("edit.html",
	"view.html", ))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func (page *Page) save() error {
	filename := "data/" + page.Title + ".txt"
	return ioutil.WriteFile(filename, page.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "data/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(writer http.ResponseWriter, request *http.Request, title string) {
	page, error := loadPage(title)
	if error != nil {
		http.Redirect(writer, request, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(writer, "view", page)
}

func editHandler(writer http.ResponseWriter, request *http.Request, title string) {
	page, error := loadPage(title)
	if error != nil {
		page = &Page{Title: title}
	}
	renderTemplate(writer, "edit", page)
}

func saveHandler(writer http.ResponseWriter, request *http.Request, title string) {
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

func makeHandler(function func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		isValidPath := validPath.FindStringSubmatch(request.URL.Path)
		if isValidPath == nil {
			http.NotFound(writer, request)
			return
		}
		function(writer, request, isValidPath[2])
	}
}

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "Welcome to bobos new website")
	})
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
