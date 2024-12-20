package session

import (
	"reflect"
	"testing"
)

func TestNewContext(t *testing.T) {
	type args struct {
		query  string
		schema string
	}
	tests := []struct {
		name string
		args args
		want *Context
	}{
		{name: "test", args: args{query: "test", schema: "test"}, want: &Context{Query: "test", Schema: "test"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewContext(tt.args.schema); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewContext() = %v, want %v", got, tt.want)
			}
		})
	}
}
