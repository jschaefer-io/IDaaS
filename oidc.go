package main

import (
	"html/template"
	"net/http"
	"sync"
)

func (s Server) AuthRequest() http.HandlerFunc {
	var init sync.Once
	var tpl *template.Template
	var tplErr error
	type templateData struct {
		Greeting string
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		init.Do(func() {
			tpl, tplErr = template.ParseFiles("./templates/index.html")
		})
		if tplErr != nil {
			http.Error(writer, tplErr.Error(), http.StatusInternalServerError)
		}
		tpl.Execute(writer, templateData{"World"})
	}
}
