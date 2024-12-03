package version

import (
	"reflect"
	"testing"
)

func TestAppInfo(t *testing.T) {
	tests := []struct {
		name string
		want BuildInfo
	}{
		{name: "test", want: BuildInfo{Version: "dev", BuildDate: "now", Codename: "capybara", GoVersion: goVersion, GoOS: goOS, GoArch: goArch}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AppInfo(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
