package item

import "github.com/zemirco/couchdb"

type Todo struct {
	couchdb.Document
	Type      string  `json:"type"`
	Text      string  `json:"text"`
	CreatedAt float64 `json:"createdAt"`
	Done      bool    `json:"done"`
}

func NewTodo(text string) *Todo {
	return &Todo{
		Type: "todo",
		Text: text,
		Done: false,
	}
}
