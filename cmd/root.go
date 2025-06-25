package cmd

import (
	"dbcli/cmd/admin"
	"dbcli/cmd/commands"
	"dbcli/internal/config"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&config.Verbose, "verbose", "v", false, "Verbose logs")
	rootCmd.PersistentFlags().StringVarP(&config.DBPath, "db", "d", "", "Path to SQLite DB file")

	rootCmd.AddCommand(admin.CmdAdmin)
	rootCmd.AddCommand(commands.CmdSql)
	rootCmd.AddCommand(commands.CmdListTables)
	rootCmd.AddCommand(commands.CmdSelect)
}

var rootCmd = &cobra.Command{Use: "dbcli"}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
