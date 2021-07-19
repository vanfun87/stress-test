package cmd

import (
	"net/http"

	"github.com/ginkgoch/stress-test/pkg/client"
	"github.com/ginkgoch/stress-test/pkg/templates"
	"github.com/spf13/cobra"
)

var (
	requestCount    int
	concurrentCount int
	requestVerb     string
	keepAlive       string
)

func init() {
	curlCmd.PersistentFlags().IntVarP(&requestCount, "requestCount", "c", 20000, "e.g 20000")
	curlCmd.PersistentFlags().IntVarP(&concurrentCount, "concurrentCount", "p", 100, "e.g 100")
	curlCmd.PersistentFlags().StringVarP(&requestVerb, "requestVerb", "v", "GET", "GET|POST|PUT|DELETE")
	curlCmd.PersistentFlags().StringVarP(&keepAlive, "keepAlive", "k", "true", "true|t|1 or false|f|0")

	rootCmd.AddCommand(curlCmd)
}

var curlCmd = &cobra.Command{
	Use:   "curl <url>",
	Short: "Curl an url",
	Long:  `Curl an url`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		s := client.NewStressClientWithConcurrentNumber(requestCount, concurrentCount)

		s.Header()
		s.Run(func() error {
			request, _ := http.NewRequest(requestVerb, args[0], nil)
			err := templates.HttpGet(request)
			return err
		})
	},
}
