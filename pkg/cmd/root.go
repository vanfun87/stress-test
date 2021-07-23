package cmd

import (
	"fmt"
	"os"

	"github.com/ginkgoch/stress-test/pkg/client/statistics"
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
	rootCmd.PersistentFlags().BoolVarP(&statistics.EnableLogger, "log", "o", false, "-o, default false")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
