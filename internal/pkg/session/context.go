// Package user contains the user context
package session

// Context represents the user context
type Context struct {
	Query  string
	Schema string
}

func NewContext(schema string) *Context {
	return &Context{Schema: schema}
}
