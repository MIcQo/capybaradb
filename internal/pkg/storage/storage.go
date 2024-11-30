package storage

// TableStorageEngine is the interface for the storage engine
type TableStorageEngine interface {
	EngineName() string
}

// SchemaStorage is the temporary storage
var SchemaStorage = map[string]*Schema{
	"public": {
		Name:   "public",
		Tables: map[string]Table{},
	},
}
