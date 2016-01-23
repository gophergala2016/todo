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
	// should be POST
	r.Methods("GET").Path("/logout").Handler(appHandler(LogoutHandler))
	log.Fatal(http.ListenAndServe(":3000", r))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) *appError {
	// check for logged in user
	session, err := RediStore.Get(r, "session")
	if err != nil {
		return InternalServerError(fmt.Errorf("get session from redistore: %v", err))
	}
	username := session.Values["username"]
	if username == nil {
		// anonymous user
		data := struct {
			Username string
		}{
			Username: "",
		}
		return indexTemplate.Execute(w, data)
	}
	// logged in user
	susername, ok := username.(string)
	if !ok {
		return InternalServerError(fmt.Errorf("parse %v to string", username))
	}
	data := struct {
		Username string
	}{
		Username: susername,
	}
	return indexTemplate.Execute(w, data)
}

func GetLoginHandler(w http.ResponseWriter, r *http.Request) *appError {
	// check for logged in user
	session, err := RediStore.Get(r, "session")
	if err != nil {
		return InternalServerError(fmt.Errorf("get session from redistore: %v", err))
	}
	username := session.Values["username"]
	if username == nil {
		// anonymous user
		data := struct {
			Username string
		}{
			Username: "",
		}
		return loginTemplate.Execute(w, data)
	}
	// logged in user
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
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

func LogoutHandler(w http.ResponseWriter, r *http.Request) *appError {
	session, err := RediStore.Get(r, "session")
	if err != nil {
		return InternalServerError(fmt.Errorf("get session from redistore: %v", err))
	}
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		return InternalServerError(fmt.Errorf("save session: %v", err))
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
