package engine

import (
	"capybaradb/internal/pkg/storage"
	"capybaradb/internal/pkg/user"

	"github.com/sirupsen/logrus"
	"vitess.io/vitess/go/vt/sqlparser"
)

// UseDatabaseStatement represents a use database statement
type UseDatabaseStatement struct {
	storage storage.Storage
}

// NewUseDatabaseStatement creates a new use database statement
func NewUseDatabaseStatement(storage storage.Storage) *UseDatabaseStatement {
	return &UseDatabaseStatement{
		storage: storage,
	}
}

// Execute executes a use database statement
func (a *UseDatabaseStatement) Execute(userContext *user.Context, s sqlparser.Statement) (StatementResult, error) {
	var v = s.(*sqlparser.Use)

	var dbName = v.DBName
	logrus.WithField("query", userContext.Query).
		WithField("schema", dbName.String()).
		Trace("switching schema")

	if dbName.IsEmpty() {
		return NewEmptyResult(), nil
	}

	if !a.storage.HasSchema(dbName.String()) {
		return NewEmptyResult(), errUnknownSchema
	}

	userContext.Schema = dbName.String()

	useCounter.WithLabelValues().Inc()

	return NewUpdateResult(1), nil
}
