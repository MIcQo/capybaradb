package engine

import (
	"capybaradb/internal/pkg/storage"
	"capybaradb/internal/pkg/user"
	"errors"

	"vitess.io/vitess/go/vt/sqlparser"
)

var defaultEngine = storage.NewRowEngine()

// CreateTableStatement creates a new create table statement
type CreateTableStatement struct {
	storage storage.Storage
}

// NewCreateTableStatement creates a new table executor
func NewCreateTableStatement(storage storage.Storage) *CreateTableStatement {
	return &CreateTableStatement{
		storage: storage,
	}
}

// Execute creates a new table
func (c *CreateTableStatement) Execute(userContext *user.Context, s sqlparser.Statement) (StatementResult, error) {
	var v = s.(*sqlparser.CreateTable)

	if v.Table.IsEmpty() {
		return NewEmptyResult(), errors.New("table name is empty")
	}

	if c.storage.HasTable(userContext.Schema, v.Table.Name.String()) {
		return NewEmptyResult(), errors.New("table already exists")

	}

	var engine = defaultEngine
	var columns = make([]storage.Column, 0)

	if err := c.storage.CreateTable(userContext.Schema, engine, v.Table.Name.String(), columns, ""); err != nil {
		return NewEmptyResult(), err
	}

	return NewUpdateResult(1), nil
}
