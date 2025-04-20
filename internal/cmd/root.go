package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var verbose bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose logs")
}

var rootCmd = &cobra.Command{Use: "dbcli"}

func Execute() {
	rootCmd.AddCommand(cmdSelect)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
