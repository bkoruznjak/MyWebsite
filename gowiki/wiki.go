package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"log"
)

type Page struct {
	Title string
	Body  []byte
}

var templatePages = template.Must(template.ParseFiles("edit.html",
	"view.html", ))

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

func viewHandler(writer http.ResponseWriter, request *http.Request) {
	title := request.URL.Path[len("/view/"):]
	page, error := loadPage(title)
	if error != nil {
		http.Redirect(writer, request, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(writer, "view", page)
}

func editHandler(writer http.ResponseWriter, request *http.Request) {
	title := request.URL.Path[len("/edit/"):]
	page, error := loadPage(title)
	if error != nil {
		page = &Page{Title: title}
	}
	renderTemplate(writer, "edit", page)
}

func saveHandler(writer http.ResponseWriter, request *http.Request) {
	title := request.URL.Path[len("/save/"):]
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
