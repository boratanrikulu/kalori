package controllers

import (
	"html/template"
	"net/http"

	"github.com/boratanrikulu/kalori/controllers/helpers"
)

func WelcomeGet(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.
		ParseFiles(helpers.GetTemplateFiles("./views/welcome/index.html")...))
	tmpl.Execute(w, nil)
}
