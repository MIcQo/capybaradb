package storage

type RowEngine struct {
	// TODO: implement row engine
}

func NewRowEngine() *RowEngine {
	return &RowEngine{}
}

func (row *RowEngine) EngineName() string {
	return "row"
}
