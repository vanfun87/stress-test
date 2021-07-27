package cmd

import (
	"net/http"
	"strings"

	"github.com/ginkgoch/stress-test/pkg/client"
	"github.com/ginkgoch/stress-test/pkg/templates"
	"github.com/spf13/cobra"
	"go.uber.org/ratelimit"
)

var (
	requestCount    int
	concurrentCount int
	requestVerb     string
	headers         []string
)

func init() {
	curlCmd.PersistentFlags().IntVarP(&requestCount, "requestCount", "c", 20000, "e.g 20000")
	curlCmd.PersistentFlags().IntVarP(&concurrentCount, "concurrentCount", "p", 100, "e.g 100")
	curlCmd.PersistentFlags().StringVarP(&requestVerb, "requestVerb", "v", "GET", "GET|POST|PUT|DELETE")
	curlCmd.PersistentFlags().StringArrayVarP(&headers, "header", "H", []string{}, "origin=eureka.com")

	rootCmd.AddCommand(curlCmd)
}

var curlCmd = &cobra.Command{
	Use:     "curl <url>",
	Short:   "Curl an url",
	Long:    `Curl an url`,
	Args:    cobra.MinimumNArgs(1),
	Example: `stress-test curl http://localhost:3000/version -c 10000 -p 100 -H origin=moblab.com -H authorization="bearer abc" -k f`,
	Run: func(cmd *cobra.Command, args []string) {
		httpClient := NewHttpClient(ParseBool(keepAlive))

		s := client.NewStressClientWithConcurrentNumber(requestCount, concurrentCount)

		var rateLimiter ratelimit.Limiter
		if limit > 0 {
			rateLimiter = ratelimit.New(limit)
		}

		s.Header()
		s.RunSingleTaskWithRateLimiter("curl", rateLimiter, func() error {
			request, _ := http.NewRequest(requestVerb, args[0], nil)

			if len(headers) > 0 {
				for _, header := range headers {
					segs := strings.Split(header, "=")
					request.Header.Set(segs[0], segs[1])
				}
			}

			err := templates.HttpGet(request, httpClient)
			return err
		})
	},
}
