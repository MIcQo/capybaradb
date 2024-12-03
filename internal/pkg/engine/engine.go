package engine

import (
	"capybaradb/internal/pkg/storage"
	"capybaradb/internal/pkg/user"
	"errors"

	"github.com/sirupsen/logrus"
	"vitess.io/vitess/go/vt/sqlparser"
)

// DefaultStorageEngine is the default storage engine
var DefaultStorageEngine = storage.NewInMemoryStorage()

var (
	errEngineNotFound   = errors.New("engine not found") // nolint
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
func ExecuteStatement(
	config *Config,
	userContext *user.Context,
	stmt sqlparser.Statement,
) (StatementResult, error) {
	var dbStorage = DefaultStorageEngine

	if config.defaultSchemaName != "" {
		if err := dbStorage.CreateSchema(config.defaultSchemaName, "default schema"); err != nil {
			panic(err)
		}
	}

	var executor Statement
	switch v := stmt.(type) {
	case *sqlparser.CreateDatabase:
		executor = NewCreateDatabaseStatement(dbStorage)
	case *sqlparser.Use:
		executor = NewUseDatabaseStatement(dbStorage)
	case *sqlparser.Select:
		executor = NewSelectStatement()
	case *sqlparser.Show:
		executor = NewShowStatement(dbStorage)
	case *sqlparser.CreateTable:
		executor = NewCreateTableStatement(dbStorage)
	default:
		logrus.Debugf("Unknown statement %+#v", v)
		return Result{}, errUnknownStatement
	}

	return executor.Execute(userContext, stmt)
}
