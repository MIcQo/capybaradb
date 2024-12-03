package storage

// Column defines a table column
type Column struct {
	Name     string
	DataType string
	NotNull  bool
	Unique   bool
}

// Row represents a single row in a table
type Row map[string]any

// Table is the storage table
type Table struct {
	Name       string   `json:"name"`
	Columns    []Column `json:"columns"`
	Rows       []Row    `json:"rows"`
	PrimaryKey string   `json:"primary_key"`
	Engine     TableStorageEngine
}

// NewTable returns a new table
func NewTable(engine TableStorageEngine, name string, columns []Column, primaryKey string) *Table {
	return &Table{
		Engine:     engine,
		Name:       name,
		Columns:    columns,
		PrimaryKey: primaryKey,
	}
}
