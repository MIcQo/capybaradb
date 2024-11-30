package engine

import (
	"capybaradb/internal/pkg/storage"
	"capybaradb/internal/pkg/user"

	"vitess.io/vitess/go/vt/sqlparser"
)

type ShowStatement struct {
}

func NewShowStatement() *ShowStatement {
	return &ShowStatement{}
}

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
		default:
			return NewEmptyResult(), nil
		}
	}

	return NewEmptyResult(), nil
}
