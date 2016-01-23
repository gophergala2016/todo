package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	indexTemplate = NewAppTemplate("index.html")
	loginTemplate = NewAppTemplate("login.html")
)

func main() {
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.Methods("GET").Path("/").Handler(appHandler(IndexHandler))
	r.Methods("GET").Path("/login").Handler(appHandler(GetLoginHandler))
	r.Methods("POST").Path("/login").Handler(appHandler(PostLoginHandler))
	log.Fatal(http.ListenAndServe(":3000", r))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) *appError {
	return indexTemplate.Execute(w, nil)
}

func GetLoginHandler(w http.ResponseWriter, r *http.Request) *appError {
	return loginTemplate.Execute(w, nil)
}

func PostLoginHandler(w http.ResponseWriter, r *http.Request) *appError {
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "john" && password == "john" {
		session, err := RediStore.Get(r, "session")
		if err != nil {
			return InternalServerError(fmt.Errorf("get session from redistore: %v", err))
		}
		session.Values["username"] = username
		if err := session.Save(r, w); err != nil {
			return InternalServerError(fmt.Errorf("save session: %v", err))
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	}
	// invalid user
	http.Redirect(w, r, r.URL.String(), http.StatusSeeOther)
	return nil
}
