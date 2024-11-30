package storage

// Schema is the storage schema
type Schema struct {
	Tables      map[string]Table
	Name        string
	Description string
}
