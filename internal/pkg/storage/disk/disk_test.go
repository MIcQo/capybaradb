package disk

import (
	"reflect"
	"testing"
)

func TestUsage(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name     string
		args     args
		wantDisk Status
	}{
		{name: "test", args: args{path: "test"}, wantDisk: Status{All: 0, Used: 0, Free: 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDisk := Usage(tt.args.path); !reflect.DeepEqual(gotDisk, tt.wantDisk) {
				t.Errorf("Usage() = %v, want %v", gotDisk, tt.wantDisk)
			}
		})
	}
}
