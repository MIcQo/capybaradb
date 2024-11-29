package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	Version   = "dev"
	BuildDate = "now"
	Codename  = "capybara"
	goVersion = runtime.Version()
	goOS      = runtime.GOOS
	goArch    = runtime.GOARCH
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the version information",
	Long:  `Displays detailed information about the build version, build date, code name, and runtime details.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version:    %s\n", Version)
		fmt.Printf("Build Date: %s\n", BuildDate)
		fmt.Printf("Code Name:  %s\n", Codename)
		fmt.Printf("Go Version: %s\n", goVersion)
		fmt.Printf("OS/Arch:    %s/%s\n", goOS, goArch)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
