package storage

import (
	"capybaradb/internal/pkg/storage/disk"
	"encoding/gob"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

const (
	dataDir = "./data"

	generalBinFile = "capybara.bin"
)

type DiskStorage struct {
	schemas []Schema
}

func NewDiskStorage() *DiskStorage {
	var schemas = make([]Schema, 0)
	var logger = logrus.WithField("file", fmt.Sprintf("%s/%s", dataDir, generalBinFile))

	if _, err := os.Stat(fmt.Sprintf("%s/%s", dataDir, generalBinFile)); err == nil {
		logger.Debug("file found, loading schemas")
		schemas, err = loadSchemaBinary(dataDir, generalBinFile)
		if err != nil {
			logrus.Fatal(err)
		}
	} else {
		logger.Debug("no file found, creating a new one with default schema")
		var newSchemas = []Schema{
			{
				Tables:      map[string]*Table{},
				Name:        "general",
				Description: "default schema",
			},
		}

		err = saveSchemaBinary(newSchemas, dataDir)
		if err != nil {
			logrus.Fatal(err)
		}

		var createSchemaDirErr = createSchemaDir(newSchemas[0].String())
		if createSchemaDirErr != nil {
			logrus.Fatal(createSchemaDirErr)
		}

		schemas = newSchemas
	}

	return &DiskStorage{
		schemas: schemas,
	}
}

func (d *DiskStorage) CreateSchema(name, description string) error {
	d.schemas = append(d.schemas, Schema{
		Tables:      map[string]*Table{},
		Name:        name,
		Description: description,
	})

	var err = saveSchemaBinary(d.schemas, dataDir)
	if err != nil {
		return err
	}

	return createSchemaDir(name)
}

func (d *DiskStorage) HasSchema(name string) bool {
	//TODO implement me
	panic("implement me")
}

func (d *DiskStorage) ListSchemas() ([][]string, error) {
	var schemas = make([][]string, 0)

	for _, schema := range d.schemas {
		schemas = append(schemas, []string{schema.String(), schema.Description})
	}

	return schemas, nil
}

func (d *DiskStorage) HasTable(schema, table string) bool {
	//TODO implement me
	panic("implement me")
}

func (d *DiskStorage) CreateTable(schema string, engine TableStorageEngine, name string, columns []Column, primaryKey string) error {
	//TODO implement me
	panic("implement me")
}

func (d *DiskStorage) ListTables(schema string) ([][]string, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DiskStorage) Free() (uint64, error) {
	var status = disk.Usage(".")
	return status.Free, nil
}

// loadSchemaBinary loads a schema from a binary file
func loadSchemaBinary(directory, schemaName string) ([]Schema, error) {
	// Open the binary schema file
	schemaFile := filepath.Join(directory, schemaName)
	file, err := os.Open(schemaFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode binary data into a Schema
	var schemas []Schema
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&schemas); err != nil {
		return nil, err
	}

	return schemas, nil
}

func saveSchemaBinary(schemas []Schema, directory string) error {
	// Ensure the directory exists
	if err := os.MkdirAll(directory, 0755); err != nil {
		return err
	}

	// Create a binary file for the schema
	schemaFile := filepath.Join(dataDir, generalBinFile)
	file, err := os.Create(schemaFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode schema to binary
	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(schemas); err != nil {
		return err
	}

	return nil
}

func createSchemaDir(schema string) error {
	return os.Mkdir(fmt.Sprintf("%s/%s", dataDir, schema), 0755)
}
