package storage

// TableStorageEngine is the interface for the storage engine
type TableStorageEngine interface {
	EngineName() string
}

type Storage interface {
	CreateSchema(name, description string) error
	HasSchema(name string) bool
	ListSchemas() ([][]string, error)

	HasTable(schema, table string) bool
	CreateTable(schema string, engine TableStorageEngine, name string, columns []Column) error
	ListTables(schema string) ([][]string, error)

	// Free returns the free space
	Free() (uint64, error)
}
