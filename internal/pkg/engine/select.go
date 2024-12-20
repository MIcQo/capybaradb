package engine

import (
	"capybaradb/internal/pkg/session"
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
func (SelectStatement) Execute(_ *session.Context, s sqlparser.Statement) (StatementResult, error) {
	var v = s.(*sqlparser.Select)
	//var table = v.GetFrom()[0].(*sqlparser.AliasedTableExpr)

	fmt.Println("")
	for _, column := range v.GetColumns() {
		switch c := column.(type) {
		case *sqlparser.AliasedExpr:
			fmt.Printf("Expr: %+#v\n", c.Expr.(*sqlparser.Variable).Name.String())
			fmt.Printf("As: %+#v\n", c.As.String())
		case *sqlparser.StarExpr:
			fmt.Printf("Expr: %+#v\n", c)
		}
	}

	fmt.Println("")

	for _, table := range v.GetFrom() {
		var t = table.(*sqlparser.AliasedTableExpr)

		fmt.Printf("Table: %+#v\n", t.TableNameString())
		fmt.Printf("Hints: %+#v\n", t.Hints)
		fmt.Printf("Alias: %+#v\n", t.As)
		fmt.Printf("Columns: %+#v\n", t.Columns)
	}

	fmt.Println("")
	fmt.Printf("Where: %+#v\n", v.Where)

	fmt.Println("")
	fmt.Printf("Limit: %+#v\n", v.Limit.Rowcount)
	fmt.Printf("Offset: %+#v\n", v.Limit.Offset)
	fmt.Println("")

	//switch v := tableExpr.(type) {
	//case *sqlparser.AliasedTableExpr:
	//	fmt.Printf("AliasedTableExpr: %v\n", v)
	//case *sqlparser.JoinTableExpr:
	//	fmt.Printf("JoinTableExpr: %v\n", v)
	//case *sqlparser.ParenTableExpr:
	//	fmt.Printf("ParenTableExpr: %v\n", v)
	////default:
	////	return errUnknownStatement
	//}

	//fmt.Printf("From: %+#v -> %d\n", from, len(from))

	fmt.Printf("Select statement: %+#v\n", s)

	selectCounter.WithLabelValues().Inc()

	return NewSelectResult(0, nil, nil), nil
}
