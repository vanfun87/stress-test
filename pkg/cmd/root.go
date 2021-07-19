package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "stress-test",
	Short: "stress-test provides a concurrent way of doing one task",
	Long:  `stress-test provides a concurrent way of doing one task with specific number, metrics will automatically printed in the terminal`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
