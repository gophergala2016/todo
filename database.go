package main

import (
	"fmt"

	"github.com/zemirco/couchdb"
	"github.com/zemirco/todo/item"
)

type Database struct {
	Client *couchdb.Client
}

// CreateUser creates per user database and saves user document to _users database
func (d Database) CreateUser(user couchdb.User) error {
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
func (d Database) SaveTodo(database string, todo *item.Todo) error {
	db := d.Client.Use(database)
	if _, err := db.Post(todo); err != nil {
		return fmt.Errorf("post todo: %v", err)
	}
	return nil
}
