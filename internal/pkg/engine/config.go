package engine

import "capybaradb/internal/pkg/storage"

// Option is a function that can be passed to NewConfig
type Option func(*Config)

// Config is the configuration for the engine
type Config struct {
	defaultSchemaName string

	storage storage.Storage
}

// WithDefaultSchema sets the default schema
func WithDefaultSchema(name string) Option {
	return func(c *Config) {
		if name == "" {
			panic("schema name is empty")
		}

		c.defaultSchemaName = name
	}
}

// WithStorage sets the storage
func WithStorage(storage storage.Storage) Option {
	return func(c *Config) {
		c.storage = storage
	}
}

// NewConfig creates a new Config for engine
func NewConfig(l ...Option) *Config {
	var c = &Config{}

	for _, v := range l {
		v(c)
	}

	return c
}
