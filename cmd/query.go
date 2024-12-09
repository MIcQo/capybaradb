// Package cmd contains application main commands
package cmd

import (
	"capybaradb/internal/pkg/config"
	"capybaradb/internal/pkg/tcp"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"net"
	"os"
	"strings"
	"time"

	"github.com/chzyer/readline"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
	queryCmd.PersistentFlags().StringP("host", "H", config.DefaultHost, "Host to connect to")
	queryCmd.PersistentFlags().StringP("user", "u", "root", "Username to connect with")
	queryCmd.PersistentFlags().StringP("password", "P", "", "Password to connect with")
	queryCmd.PersistentFlags().UintP("port", "p", config.DefaultDatabasePort, "Port to connect to")

	queryCmd.PersistentFlags().BoolP("interactive", "i", false, "Enable interactive mode toggle")

	rootCmd.AddCommand(queryCmd)
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

	var result tcp.Packet
	var decoder = gob.NewDecoder(strings.NewReader(string(out)))

	if err := decoder.Decode(&result); err != nil {
		logrus.WithError(err).Debug("Failed to decode response from database")
		return
	}

	if len(result.Rows) > 0 {
		var t = table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		var headers = table.Row{}
		for _, header := range result.Columns {
			headers = append(headers, header)
		}

		t.AppendHeader(headers)

		for _, row := range result.Rows {
			var r = table.Row{}
			for _, cell := range row {
				r = append(r, cell)
			}

			t.AppendRow(r)
		}

		t.Render()
	} else {
		fmt.Printf(
			"Affected rows: %d\n",
			result.AffectedRows,
		)
	}

	fmt.Printf("Query took: %s\n", time.Since(n))
}
