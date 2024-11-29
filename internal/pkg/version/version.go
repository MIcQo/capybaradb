package version

import "runtime"

var (
	Version   = "dev"
	BuildDate = "now"
	Codename  = "capybara"
	goVersion = runtime.Version()
	goOS      = runtime.GOOS
	goArch    = runtime.GOARCH
)

type BuildInfo struct {
	Version   string
	BuildDate string
	Codename  string
	GoVersion string
	GoOS      string
	GoArch    string
}

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
