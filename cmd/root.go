package cmd

import (
	"capybaradb/internal/pkg/version"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "capybaradb",
	Short: "CapybaraDB is a scalable database for hybrid workloads",
	Long: `CapybaraDB is a high-performance, ACID-compliant database 
designed for both transactional and analytical workloads. 
It supports clustering and hybrid storage models (row and column-oriented).`,
	// Funkcia vykonaná, keď nie je uvedený žiadny podpríkaz
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to CapybaraDB!")
		fmt.Println("Use the --help flag to see available commands.")
	},

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Flags().GetBool("debug")
		if debug {
			logrus.Info("Debug mode enabled")
			logrus.SetLevel(logrus.DebugLevel)
		}

		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			DisableColors:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			FieldMap:        logrus.FieldMap{"version": version.Version},
		})
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var Debug bool

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringP("config", "c", "", "Path to the configuration file")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.capybaradb.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
