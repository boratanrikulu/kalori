package controllers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/boratanrikulu/kalori/foods"
)

const (
	MB = 1 << 20
)

func RecognizePost(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 * MB)
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	var b bytes.Buffer
	io.Copy(&b, file)

	err = checkMime(b.Bytes())
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	result, probability := foods.Recognize(&b, "png")
	output := fmt.Sprint(result, " - ", probability)
	if probability <= 0.5 {
		output = "Could not found."
	}

	fmt.Fprintln(w, output)
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
