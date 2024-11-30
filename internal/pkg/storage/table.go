package storage

// Table is the storage table
type Table struct {
	Name    string
	Engine  TableStorageEngine
	Columns []Column
}

// Column is the interface for the column engine
type Column struct {
	Name string
}
