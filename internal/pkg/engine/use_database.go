package engine

import (
	"capybaradb/internal/pkg/storage"
	"capybaradb/internal/pkg/user"

	"github.com/sirupsen/logrus"
	"vitess.io/vitess/go/vt/sqlparser"
)

type UseDatabaseStatement struct {
}

func NewUseDatabaseStatement() *UseDatabaseStatement {
	return &UseDatabaseStatement{}
}

func (UseDatabaseStatement) Execute(userContext *user.Context, s sqlparser.Statement) (StatementResult, error) {
	var v = s.(*sqlparser.Use)

	var dbName = v.DBName
	logrus.WithField("query", userContext.Query).
		WithField("schema", dbName.String()).
		Trace("switching schema")

	if dbName.IsEmpty() {
		return NewEmptyResult(), nil
	}

	if _, ok := storage.SchemaStorage[dbName.String()]; !ok {
		return NewEmptyResult(), errUnknownSchema
	}

	userContext.Schema = dbName.String()

	useCounter.WithLabelValues().Inc()

	return NewUpdateResult(1), nil
}
