package cmd

import (
	"capybaradb/internal/pkg/version"
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the version information",
	Long:  `Displays detailed information about the build version, build date, code name, and runtime details.`,
	Run: func(_ *cobra.Command, _ []string) {
		var info = version.AppInfo()

		fmt.Printf("Version:    %s\n", info.Version)
		fmt.Printf("Build Date: %s\n", info.BuildDate)
		fmt.Printf("Code Name:  %s\n", info.Codename)
		fmt.Printf("Go Version: %s\n", info.GoVersion)
		fmt.Printf("OS/Arch:    %s/%s\n", info.GoOS, info.GoArch)
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
