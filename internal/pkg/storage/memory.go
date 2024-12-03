package storage

import (
	"math"
)

// InMemoryStorage is the in memory storage for testing purpose
type InMemoryStorage struct {
	inMemory map[string]*Schema
}

func (m *InMemoryStorage) CreateTable(schema string, engine TableStorageEngine, name string, columns []Column) error {
	m.inMemory[schema].Tables[name] = Table{
		Name:    name,
		Engine:  engine,
		Columns: columns,
	}
	return nil
}

func (m *InMemoryStorage) ListTables(schema string) ([][]string, error) {
	var tables = make([][]string, 0)
	for _, table := range m.inMemory[schema].Tables {
		tables = append(tables, []string{table.Name})
	}

	return tables, nil
}

func (m *InMemoryStorage) HasTable(schema, table string) bool {
	_, ok := m.inMemory[schema].Tables[table]
	return ok
}

func (m *InMemoryStorage) ListSchemas() ([][]string, error) {
	var dbs = make([][]string, 0)

	for _, db := range m.inMemory {
		dbs = append(dbs, []string{db.String(), db.Description})
	}

	return dbs, nil
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		inMemory: map[string]*Schema{},
	}
}

func (m *InMemoryStorage) Free() (uint64, error) {
	return math.MaxUint64, nil
}

func (m *InMemoryStorage) CreateSchema(name, description string) error {
	m.inMemory[name] = &Schema{Name: name, Description: description}
	return nil
}

func (m *InMemoryStorage) HasSchema(name string) bool {
	_, ok := m.inMemory[name]
	return ok
}
