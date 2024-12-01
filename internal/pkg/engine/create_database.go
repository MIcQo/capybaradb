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
	storage storage.Storage
}

// NewCreateDatabaseStatement creates a new create database statement
func NewCreateDatabaseStatement(storage storage.Storage) *CreateDatabaseStatement {
	return &CreateDatabaseStatement{
		storage: storage,
	}
}

// Execute creates a new database
func (a *CreateDatabaseStatement) Execute(userContext *user.Context, s sqlparser.Statement) (StatementResult, error) {
	var v = s.(*sqlparser.CreateDatabase)

	var dbName = v.DBName
	logrus.WithField("query", userContext.Query).
		WithField("schema", dbName.String()).
		Trace("creating schema")

	if v.IfNotExists && a.storage.HasSchema(dbName.String()) {
		return NewEmptyResult(), nil
	}

	var err = a.storage.CreateSchema(dbName.String(), "")
	if err != nil {
		return NewEmptyResult(), err
	}

	createDatabaseCounter.WithLabelValues().Inc()

	return Result{affectedRows: 1}, nil
}
