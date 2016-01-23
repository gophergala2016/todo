package item

// Todo defines todo item.
// Whole package is reused for mobile app via gomobile.
// Therefore only basic types work.
type Todo struct {
	// cannot embed couchdb.Document directly because gomobile complains
	ID        string  `json:"_id,omitempty"`
	Rev       string  `json:"_rev,omitempty"`
	Type      string  `json:"type"`
	Text      string  `json:"text"`
	CreatedAt float64 `json:"createdAt"`
	Done      bool    `json:"done"`
}

// NewTodo returns new todo item.
// Constructor must exist and must return pointer to work on ios.
func NewTodo(text string) *Todo {
	return &Todo{
		Type: "todo",
		Text: text,
		Done: false,
	}
}

// let Todo implement couchdb.CouchDoc interface

// GetID returns document ID.
func (t *Todo) GetID() string {
	return t.ID
}

// GetRev returns document revision.
func (t *Todo) GetRev() string {
	return t.Rev
}
