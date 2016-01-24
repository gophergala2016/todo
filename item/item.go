package item

// Todo describes todo item.
type Todo struct {
	ID        string  `json:"_id,omitempty"`
	Rev       string  `json:"_rev,omitempty"`
	Type      string  `json:"type"`
	Text      string  `json:"text"`
	CreatedAt float64 `json:"createdAt"`
	Done      bool    `json:"done"`
}

// NewTodo returns new todo item.
func NewTodo(text string) *Todo {
	return &Todo{
		Type: "todo",
		Text: text,
		Done: false,
	}
}

// implement couchdb.CouchDoc interface

// GetID returns document ID.
func (t *Todo) GetID() string {
	return t.ID
}

// GetRev returns document revision.
func (t *Todo) GetRev() string {
	return t.Rev
}
