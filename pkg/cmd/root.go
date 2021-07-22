package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var keepAlive string

var rootCmd = &cobra.Command{
	Use:   "stress-test",
	Short: "stress-test provides a concurrent way of doing one task",
	Long:  `stress-test provides a concurrent way of doing one task with specific number, metrics will automatically printed in the terminal`,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&keepAlive, "keepAlive", "k", "true", "true|t|1 or false|f|0")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
