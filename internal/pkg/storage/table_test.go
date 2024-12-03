package storage

import (
	"reflect"
	"testing"
)

func TestNewTable(t *testing.T) {
	var engine = NewRowEngine()

	type args struct {
		engine     TableStorageEngine
		name       string
		columns    []Column
		primaryKey string
	}
	tests := []struct {
		name string
		args args
		want *Table
	}{
		{name: "test", args: args{engine: engine, name: "test", columns: []Column{}, primaryKey: "test"}, want: &Table{Name: "test", Columns: []Column{}, PrimaryKey: "test", Engine: engine}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTable(tt.args.engine, tt.args.name, tt.args.columns, tt.args.primaryKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTable() = %v, want %v", got, tt.want)
			}
		})
	}
}
