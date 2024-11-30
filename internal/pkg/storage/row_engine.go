package storage

// RowEngine is the row engine
type RowEngine struct {
	// TODO: implement row engine
}

// NewRowEngine returns a new row engine
func NewRowEngine() *RowEngine {
	return &RowEngine{}
}

// EngineName returns the engine name
func (row *RowEngine) EngineName() string {
	return "row"
}
