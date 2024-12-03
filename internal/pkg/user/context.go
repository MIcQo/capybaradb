// Package user contains the user context
package user

// Context represents the user context
type Context struct {
	Query  string
	Schema string
}

func NewContext(query string, schema string) *Context {
	return &Context{Query: query, Schema: schema}
}
