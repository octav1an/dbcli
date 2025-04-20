package data

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var dbPath string

var CmdData = &cobra.Command{
	Use:   "data",
	Short: "Groups all db query operations",
}

func init() {
	CmdData.PersistentFlags().StringVarP(&dbPath, "db", "d", "", "Path to SQLite DB file")
	CmdData.MarkPersistentFlagRequired("db")

	CmdData.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := validateDb(dbPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return nil
	}
}
