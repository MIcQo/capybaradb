package engine

import (
	"capybaradb/internal/pkg/session"
	"capybaradb/internal/pkg/storage"

	"vitess.io/vitess/go/vt/sqlparser"
)

// ShowStatement represents a show statement
type ShowStatement struct {
	storage storage.Storage
}

// NewShowStatement creates a new show statement
func NewShowStatement(storage storage.Storage) *ShowStatement {
	return &ShowStatement{
		storage: storage,
	}
}

// Execute executes a show statement
func (a *ShowStatement) Execute(userContext *session.Context, s sqlparser.Statement) (StatementResult, error) {
	var v = s.(*sqlparser.Show)

	showCounter.WithLabelValues().Inc()

	switch t := v.Internal.(type) {
	case *sqlparser.ShowBasic:
		switch t.Command {
		case sqlparser.Database:
			var databases, err = a.storage.ListSchemas()
			if err != nil {
				return NewEmptyResult(), err
			}

			return NewSelectResult(len(databases), []string{"Database", "Description"}, databases), nil
		case sqlparser.Table:
			var tables, err = a.storage.ListTables(userContext.Schema)
			if err != nil {
				return NewEmptyResult(), err
			}

			return NewSelectResult(len(tables), []string{"Table"}, tables), nil
		default:
			return NewEmptyResult(), nil
		}
	}

	return NewEmptyResult(), nil
}
