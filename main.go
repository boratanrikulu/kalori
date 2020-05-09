package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/boratanrikulu/kalori/controllers"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", controllers.WelcomeGet).Methods("GET")
	r.HandleFunc("/", controllers.RecognizeGet).Methods("GET")
	r.PathPrefix("/assets").Handler(http.StripPrefix("/assets/",
		http.FileServer(http.Dir("./assets/"))))

	http.ListenAndServe(":"+os.Getenv("PORT"), r)
}
