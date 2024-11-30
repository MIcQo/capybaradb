package engine

import (
	"capybaradb/internal/pkg/user"
	"errors"

	"github.com/sirupsen/logrus"
	"vitess.io/vitess/go/vt/sqlparser"
)

var (
	errEngineNotFound   = errors.New("engine not found")
	errUnknownStatement = errors.New("unknown statement")
	errUnknownSchema    = errors.New("unknown schema")
)

// Statement - SQL statement executor
type Statement interface {
	Execute(userContext *user.Context, s sqlparser.Statement) (StatementResult, error)
}

// StatementResult - SQL statement result
type StatementResult interface {
	AffectedRows() int
	LastInsertId() int
	Columns() []string
	Rows() [][]string
}

// ExecuteStatement executes a SQL statement through tables and engines
func ExecuteStatement(userContext *user.Context, stmt sqlparser.Statement) (StatementResult, error) {
	var executor Statement
	switch v := stmt.(type) {
	case *sqlparser.CreateDatabase:
		executor = NewCreateDatabaseStatement()
	case *sqlparser.Use:
		executor = NewUseDatabaseStatement()
	case *sqlparser.Select:
		executor = NewSelectStatement()
	case *sqlparser.Show:
		executor = NewShowStatement()
	case *sqlparser.CreateTable:
		executor = NewCreateTableStatement()
	default:
		logrus.Debugf("Unknown statement %+#v", v)
		return Result{}, errUnknownStatement
	}

	return executor.Execute(userContext, stmt)
}
