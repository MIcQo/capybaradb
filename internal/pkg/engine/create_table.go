package engine

import (
	"capybaradb/internal/pkg/storage"
	"capybaradb/internal/pkg/user"
	"errors"

	"vitess.io/vitess/go/vt/sqlparser"
)

var defaultEngine = storage.NewRowEngine()

// CreateTableStatement creates a new create table statement
type CreateTableStatement struct{}

// NewCreateTableStatement creates a new create table statement
func NewCreateTableStatement() *CreateTableStatement {
	return &CreateTableStatement{}
}

// Execute creates a new table
func (c *CreateTableStatement) Execute(userContext *user.Context, s sqlparser.Statement) (StatementResult, error) {
	var v = s.(*sqlparser.CreateTable)

	if v.Table.IsEmpty() {
		return NewEmptyResult(), errors.New("table name is empty")
	}

	if _, ok := storage.SchemaStorage[userContext.Schema].Tables[v.Table.Name.String()]; ok {
		return NewEmptyResult(), errors.New("table already exists")
	}

	var engine = defaultEngine

	storage.SchemaStorage[userContext.Schema].Tables[v.Table.Name.String()] = storage.Table{
		Name: v.Table.Name.String(), Engine: engine, Columns: make([]storage.Column, 0),
	}

	return NewUpdateResult(1), nil
}
