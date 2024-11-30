package engine

import (
	"capybaradb/internal/pkg/user"
	"fmt"

	"vitess.io/vitess/go/vt/sqlparser"
)

// SelectStatement represents a select statement
type SelectStatement struct {
}

// NewSelectStatement creates a new select statement
func NewSelectStatement() *SelectStatement {
	return &SelectStatement{}
}

// Execute executes a select statement
func (SelectStatement) Execute(_ *user.Context, s sqlparser.Statement) (StatementResult, error) {
	var v = s.(*sqlparser.Select)
	var table = v.GetFrom()[0].(*sqlparser.AliasedTableExpr)

	fmt.Printf("Table: %+#v\n", table)

	//var tableExpr = from[0]
	//switch v := tableExpr.(type) {
	//case *sqlparser.AliasedTableExpr:
	//	fmt.Printf("AliasedTableExpr: %v\n", v)
	//case *sqlparser.JoinTableExpr:
	//	fmt.Printf("JoinTableExpr: %v\n", v)
	//case *sqlparser.ParenTableExpr:
	//	fmt.Printf("ParenTableExpr: %v\n", v)
	//default:
	//	return errUnknownStatement
	//}

	//fmt.Printf("From: %+#v -> %d\n", from, len(from))

	fmt.Printf("Select statement: %+#v\n", s)

	selectCounter.WithLabelValues().Inc()

	return NewSelectResult(0, nil, nil), nil
}
