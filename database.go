package main

import (
	"encoding/json"
	"fmt"

	"github.com/segmentio/pointer"
	"github.com/zemirco/couchdb"
	"github.com/zemirco/todo/item"
)

type database struct {
	Client *couchdb.Client
}

const (
	// January 1st 2013 01:00:00
	past = "1356998400"
	// December 24th 3600 12:00:00
	future = "32503593600"
)

// CreateUser creates per user database and saves user document to _users database
func (d database) CreateUser(user couchdb.User) error {
	// create per user database
	if _, err := d.Client.Create(user.Name); err != nil {
		return fmt.Errorf("create database: %v", err)
	}
	// create views in per user database
	design := couchdb.DesignDocument{
		Document: couchdb.Document{
			ID: "_design/todo",
		},
		Language: "javascript",
		Views: map[string]couchdb.DesignDocumentView{
			"byCreatedAt": {
				Map: `
function(doc) {
	if (doc.type === "todo") {
		emit(doc.createdAt);
	}
}
				`,
			},
		},
	}
	db := d.Client.Use(user.Name)
	if _, err := db.Put(&design); err != nil {
		return fmt.Errorf("put design document: %v", err)
	}
	// create user in _users database
	db = d.Client.Use("_users")
	if _, err := db.Put(&user); err != nil {
		return fmt.Errorf("put user: %v", err)
	}
	return nil
}

// SaveTodo saves todo item to CouchDB
func (d database) SaveTodo(database string, todo *item.Todo) error {
	db := d.Client.Use(database)
	if _, err := db.Post(todo); err != nil {
		return fmt.Errorf("post todo: %v", err)
	}
	return nil
}

// GetTodos gets all todos from per user database
func (d database) GetTodos(database string) ([]item.Todo, error) {
	db := d.Client.Use(database)
	view := db.View("todo")
	params := couchdb.QueryParameters{
		StartKey:    pointer.String(past),
		EndKey:      pointer.String(future),
		IncludeDocs: pointer.Bool(true),
	}
	res, err := view.Get("byCreatedAt", params)
	if err != nil {
		return nil, fmt.Errorf("get view byCreatedAt: %v", err)
	}
	docs := make([]interface{}, len(res.Rows))
	for index, row := range res.Rows {
		docs[index] = row.Doc
	}
	todos := make([]item.Todo, len(res.Rows))
	b, err := json.Marshal(docs)
	if err != nil {
		return nil, fmt.Errorf("json marshal: %v", err)
	}
	return todos, json.Unmarshal(b, &todos)
}

func (d database) GetTodoByID(database, id string) (item.Todo, error) {
	db := d.Client.Use(database)
	var t item.Todo
	return t, db.Get(&t, id)
}

func (d database) UpdateTodo(database string, todo item.Todo) error {
	db := d.Client.Use(database)
	_, err := db.Put(&todo)
	if err != nil {
		return fmt.Errorf("put todo: %v", err)
	}
	return nil
}

func (d database) DeleteTodoByID(database, id string) error {
	db := d.Client.Use(database)
	// get document first to retrieve current revision
	doc, err := d.GetTodoByID(database, id)
	if err != nil {
		return fmt.Errorf("get todo by id: %v", err)
	}
	// delete document by id and revision
	if _, err = db.Delete(&doc); err != nil {
		return fmt.Errorf("delete document: %v", err)
	}
	return nil
}
