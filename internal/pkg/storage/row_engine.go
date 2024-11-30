package storage

type ColumnEngine struct {
	// TODO: implement column engine
}

func NewColumnEngine() *ColumnEngine {
	return &ColumnEngine{}
}

func (column *ColumnEngine) EngineName() string {
	return "row"
}
