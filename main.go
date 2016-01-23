package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	indexTemplate = NewAppTemplate("index.html")
)

func main() {
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.Methods("GET").Path("/").Handler(appHandler(IndexHandler))
	log.Fatal(http.ListenAndServe(":3000", r))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) *appError {
	return indexTemplate.Execute(w, nil)
}
