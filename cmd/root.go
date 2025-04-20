package cmd

import (
	"dbcli/cmd/admin"
	"dbcli/cmd/data"
	"dbcli/internal/config"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&config.Verbose, "verbose", "v", false, "Verbose logs")
	rootCmd.AddCommand(data.CmdData)
	rootCmd.AddCommand(admin.CmdAdmin)
}

var rootCmd = &cobra.Command{Use: "dbcli"}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
