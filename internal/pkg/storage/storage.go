package storage

// TableStorageEngine is the interface for the storage engine
type TableStorageEngine interface {
	EngineName() string
}
