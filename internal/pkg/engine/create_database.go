// Package engine is where engine sits
package engine

import (
	"capybaradb/internal/pkg/storage"
	"capybaradb/internal/pkg/user"

	"github.com/sirupsen/logrus"
	"vitess.io/vitess/go/vt/sqlparser"
)

// CreateDatabaseStatement represents a create database statement
type CreateDatabaseStatement struct {
}

// NewCreateDatabaseStatement creates a new create database statement
func NewCreateDatabaseStatement() *CreateDatabaseStatement {
	return &CreateDatabaseStatement{}
}

// Execute creates a new database
func (CreateDatabaseStatement) Execute(userContext *user.Context, s sqlparser.Statement) (StatementResult, error) {
	var v = s.(*sqlparser.CreateDatabase)

	var dbName = v.DBName
	logrus.WithField("query", userContext.Query).
		WithField("schema", dbName.String()).
		Trace("creating schema")

	if _, ok := storage.SchemaStorage[dbName.String()]; ok {
		return NewEmptyResult(), nil
	}

	storage.SchemaStorage[dbName.String()] = &storage.Schema{Name: dbName.String(), Tables: map[string]storage.Table{}}

	createDatabaseCounter.WithLabelValues().Inc()

	return Result{affectedRows: 1}, nil
}
