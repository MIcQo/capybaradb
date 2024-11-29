// Package version contains information about the application
package version

import "runtime"

var (
	// Version is the current version of the application
	Version = "dev"

	// BuildDate is the date the application was built
	BuildDate = "now"

	// Codename is the codename of the application
	Codename = "capybara"

	// goVersion is the version of Go used to build the application
	goVersion = runtime.Version()

	// goOS is the operating system used to build the application
	goOS = runtime.GOOS

	// goArch is the architecture used to build the application
	goArch = runtime.GOARCH
)

// BuildInfo contains information about the application
type BuildInfo struct {
	Version   string
	BuildDate string
	Codename  string
	GoVersion string
	GoOS      string
	GoArch    string
}

// AppInfo returns information about the application
func AppInfo() BuildInfo {
	return BuildInfo{
		Version:   Version,
		BuildDate: BuildDate,
		Codename:  Codename,
		GoVersion: goVersion,
		GoOS:      goOS,
		GoArch:    goArch,
	}
}
