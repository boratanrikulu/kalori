package controllers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/boratanrikulu/kalori/controllers/helpers"
	"github.com/boratanrikulu/kalori/foods"
)

type ResultPageData struct {
	Food struct {
		Name    string
		Calorie string
		Picture string
	}
}

type ErrorPageData struct {
	Error string
}

func RecognizePost(w http.ResponseWriter, r *http.Request) {
	errorPage := template.Must(template.
		ParseFiles(helpers.GetTemplateFiles("./views/result/error.html")...))
	indexPage := template.Must(template.
		ParseFiles(helpers.GetTemplateFiles("./views/result/index.html")...))

	r.ParseMultipartForm(1 << 20) // 10 MB size limit
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		e := ErrorPageData{
			Error: err.Error(),
		}
		errorPage.Execute(w, e)
		return
	}
	defer file.Close()

	var b bytes.Buffer
	io.Copy(&b, file)

	err = checkMime(b.Bytes())
	if err != nil {
		e := ErrorPageData{
			Error: err.Error(),
		}
		errorPage.Execute(w, e)
		return
	}

	name, calorie, err := foods.Recognize(&b)

	if err != nil {
		e := ErrorPageData{
			Error: "We could not predict your food. Try an other!",
		}
		errorPage.Execute(w, e)
		return
	}

	pageData := ResultPageData{
		Food: struct {
			Name    string
			Calorie string
			Picture string
		}{
			Name:    name,
			Calorie: calorie,
			Picture: base64.StdEncoding.EncodeToString(b.Bytes()),
		},
	}

	indexPage.Execute(w, pageData)
}

func checkMime(b []byte) error {
	allowed := []string{"image/jpeg", "image/jpg", "image/gif", "image/png", "image/bmp"}
	mime := http.DetectContentType(b)
	if contains(allowed, mime) {
		return nil
	}
	return fmt.Errorf("Not allowed format: %v\nPlease upload one of the this formats: %v",
		mime, strings.Join(allowed, ", "))
}

// contains tells whether a contains x.
func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
