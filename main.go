package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/zemirco/couchdb"
	"github.com/zemirco/todo/item"
)

var (
	indexTemplate = NewAppTemplate("index.html")
	loginTemplate = NewAppTemplate("login.html")
)

func init() {
	// init logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// add dummy user to database
	john := couchdb.NewUser("john", "john", []string{})
	if err := db.CreateUser(john); err != nil {
		log.Printf("user might already exist: %v", err)
	}
}

func main() {
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.Methods("GET").Path("/").Handler(appHandler(IndexHandler))
	r.Methods("GET").Path("/login").Handler(appHandler(GetLoginHandler))
	r.Methods("POST").Path("/login").Handler(appHandler(PostLoginHandler))
	r.Methods("POST").Path("/create").Handler(appHandler(CreateHandler))
	// should be POST
	r.Methods("GET").Path("/logout").Handler(appHandler(LogoutHandler))
	r.Methods("POST").Path("/{id}/delete").Handler(appHandler(DeleteHandler))
	r.Methods("POST").Path("/{id}/done").Handler(appHandler(DoneHandler))
	r.Methods("POST").Path("/{id}/undone").Handler(appHandler(UndoneHandler))
	// add logging
	rWithLogging := handlers.LoggingHandler(os.Stdout, r)
	log.Fatal(http.ListenAndServe(":3000", rWithLogging))
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
	// get todos from database
	todos, err := db.GetTodos(susername)
	if err != nil {
		return InternalServerError(fmt.Errorf("get todos: %v", err))
	}
	data := struct {
		Username string
		Todos    []item.Todo
	}{
		Username: susername,
		Todos:    todos,
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

func CreateHandler(w http.ResponseWriter, r *http.Request) *appError {
	text := r.FormValue("text")
	todo := item.NewTodo(text)
	// createdAt has to be set manually here as ios doesn't understand type time yet
	todo.CreatedAt = float64(time.Now().UTC().Unix())
	// get username from session
	session, err := RediStore.Get(r, "session")
	if err != nil {
		return InternalServerError(fmt.Errorf("get session from redistore: %v", err))
	}
	username := session.Values["username"].(string)
	if err := db.SaveTodo(username, todo); err != nil {
		return InternalServerError(fmt.Errorf("save todo: %v", err))
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) *appError {
	vars := mux.Vars(r)
	id := vars["id"]
	// get username from session
	session, err := RediStore.Get(r, "session")
	if err != nil {
		return InternalServerError(fmt.Errorf("get session from redistore: %v", err))
	}
	username := session.Values["username"].(string)
	if err := db.DeleteTodoByID(username, id); err != nil {
		return InternalServerError(fmt.Errorf("delete todo by id: %v", err))
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func DoneHandler(w http.ResponseWriter, r *http.Request) *appError {
	vars := mux.Vars(r)
	id := vars["id"]
	// get username from session
	session, err := RediStore.Get(r, "session")
	if err != nil {
		return InternalServerError(fmt.Errorf("get session from redistore: %v", err))
	}
	username := session.Values["username"].(string)
	t, err := db.GetTodoByID(username, id)
	if err != nil {
		return InternalServerError(fmt.Errorf("get todo by id: %v", err))
	}
	t.Done = true
	if err := db.UpdateTodo(username, t); err != nil {
		return InternalServerError(fmt.Errorf("update todo: %v", err))
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func UndoneHandler(w http.ResponseWriter, r *http.Request) *appError {
	vars := mux.Vars(r)
	id := vars["id"]
	// get username from session
	session, err := RediStore.Get(r, "session")
	if err != nil {
		return InternalServerError(fmt.Errorf("get session from redistore: %v", err))
	}
	username := session.Values["username"].(string)
	t, err := db.GetTodoByID(username, id)
	if err != nil {
		return InternalServerError(fmt.Errorf("get todo by id: %v", err))
	}
	t.Done = false
	if err := db.UpdateTodo(username, t); err != nil {
		return InternalServerError(fmt.Errorf("update todo: %v", err))
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
