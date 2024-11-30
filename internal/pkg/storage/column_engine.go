// Package storage is used to manage the storage
package storage

// ColumnEngine is the column engine
type ColumnEngine struct {
	// TODO: implement column engine
}

// NewColumnEngine returns a new column engine
func NewColumnEngine() *ColumnEngine {
	return &ColumnEngine{}
}

// EngineName returns the engine name
func (column *ColumnEngine) EngineName() string {
	return "row"
}
