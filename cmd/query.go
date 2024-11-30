// Package cmd contains application main commands
package cmd

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/chzyer/readline"
	"net"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	defaultHost = "127.0.0.1"
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query [SQL statement]",
	Short: "Run a SQL-like query against the database",
	//Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var host, _ = cmd.Flags().GetString("host")
		var port, _ = cmd.Flags().GetUint("port")
		var interactive, _ = cmd.Flags().GetBool("interactive")

		var srv, err = net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
		if err != nil {
			logrus.WithError(err).Debug("Failed to connect to database")
			return
		}

		defer func(srv net.Conn) {
			_ = srv.Close()
		}(srv)

		if interactive {
			rl, err := readline.NewEx(&readline.Config{
				Prompt:            "\033[31m>>\033[0m ",
				HistoryFile:       "/tmp/readline.tmp",
				InterruptPrompt:   "^C",
				EOFPrompt:         "exit",
				HistorySearchFold: true,
			})
			if err != nil {
				logrus.Fatalf("Error initializing readline: %v", err)
			}
			defer func(rl *readline.Instance) {
				_ = rl.Close()
			}(rl)

			for {
				line, err := rl.Readline()
				if errors.Is(err, readline.ErrInterrupt) {
					if len(line) == 0 {
						break
					}
					continue
				}

				line = strings.TrimSpace(line)

				var input = strings.Trim(line, "\n")
				switch input {
				case "exit", "quit":
					fmt.Println("Goodbye")
					return
				default:
					sendAndReadResponse(srv, input)
				}
			}
		} else {
			query := args[0]
			fmt.Printf("Executing query: %s\n", query)
			sendAndReadResponse(srv, query)
		}

		//var parse, parseErr = sqlparser.New(sqlparser.Options{})
		//if parseErr != nil {
		//	logrus.WithError(parseErr).Debug("Failed to parse query")
		//}
		//
		//var tree, treeErr = parse.Parse(query)
		//if treeErr != nil {
		//	logrus.WithError(treeErr).Debug("Failed to parse query")
		//}
		//
		//fmt.Printf("Parse Tree: %+#v\n", tree)

		//var stmt, err = compiler.PrepareStatement(query)
		//if err != nil {
		//	logrus.WithError(err).Debug("Failed to prepare statement")
		//}
		//
		//fmt.Printf("Statement: %+#v\n", stmt)

	},
}

func init() {
	queryCmd.PersistentFlags().String("host", defaultHost, "Host to connect to")
	queryCmd.PersistentFlags().Uint("port", defaultDatabasePort, "Port to connect to")
	queryCmd.PersistentFlags().BoolP("interactive", "i", false, "Enable interactive mode")

	rootCmd.AddCommand(queryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// queryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// queryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func sendAndReadResponse(srv net.Conn, query string) {
	var n = time.Now()
	var bin = hex.EncodeToString([]byte(query))
	var _, writeErr = srv.Write([]byte(bin + "\n"))
	if writeErr != nil {
		logrus.WithError(writeErr).Debug("Failed to write query to database")
		return
	}

	var out = make([]byte, 1024)
	var _, readErr = srv.Read(out)
	if readErr != nil {
		logrus.WithError(readErr).Debug("Failed to read response from database")
		return
	}

	if !strings.HasSuffix(string(out), "\n") {
		out = []byte(string(out) + "\n")
	}

	fmt.Printf("%s", out)
	fmt.Printf("Query took: %s\n", time.Since(n))
}
