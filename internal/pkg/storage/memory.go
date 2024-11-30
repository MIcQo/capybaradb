package storage

import "capybaradb/internal/pkg/config"

// SchemaStorage is the temporary storage
var SchemaStorage = map[string]*Schema{
	config.DefaultSchema: {
		Name:        config.DefaultSchema,
		Tables:      map[string]Table{},
		Description: "default schema",
	},
}
