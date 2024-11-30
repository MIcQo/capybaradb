package storage

// Schema is the storage schema
type Schema struct {
	Name   string
	Tables map[string]Table
}
