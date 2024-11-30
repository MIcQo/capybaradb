package engine

import (
	"capybaradb/internal/pkg/storage"
	"capybaradb/internal/pkg/user"

	"vitess.io/vitess/go/vt/sqlparser"
)

// ShowStatement represents a show statement
type ShowStatement struct {
}

// NewShowStatement creates a new show statement
func NewShowStatement() *ShowStatement {
	return &ShowStatement{}
}

// Execute executes a show statement
func (ShowStatement) Execute(userContext *user.Context, s sqlparser.Statement) (StatementResult, error) {
	var v = s.(*sqlparser.Show)

	showCounter.WithLabelValues().Inc()

	switch t := v.Internal.(type) {
	case *sqlparser.ShowBasic:
		switch t.Command {
		case sqlparser.Database:
			var databases = make([][]string, 0)
			for name := range storage.SchemaStorage {
				databases = append(databases, []string{name})
			}

			return NewSelectResult(len(databases), []string{"Database"}, databases), nil
		case sqlparser.Table:
			var tables = make([][]string, 0)
			for _, table := range storage.SchemaStorage[userContext.Schema].Tables {
				tables = append(tables, []string{table.Name})
			}

			return NewSelectResult(len(tables), []string{"Table"}, tables), nil
		default:
			return NewEmptyResult(), nil
		}
	}

	return NewEmptyResult(), nil
}
