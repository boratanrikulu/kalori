package controllers

import (
	"html/template"
	"net/http"
)

func WelcomeGet(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./views/welcome/index.html"))
	tmpl.Execute(w, nil)
}
