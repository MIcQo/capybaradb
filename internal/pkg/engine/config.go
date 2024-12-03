package engine

// Option is a function that can be passed to NewConfig
type Option func(*Config)

// Config is the configuration for the engine
type Config struct {
	defaultSchemaName string
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

// NewConfig creates a new Config for engine
func NewConfig(l ...Option) *Config {
	var c = &Config{}

	for _, v := range l {
		v(c)
	}

	return c
}
