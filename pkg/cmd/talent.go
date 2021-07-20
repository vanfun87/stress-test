package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync/atomic"

	"github.com/ginkgoch/stress-test/pkg/client"
	"github.com/ginkgoch/stress-test/pkg/talent"
	"github.com/spf13/cobra"
	"go.uber.org/ratelimit"
)

var (
	filepath string
	limit    int
)

func init() {
	toCmd.PersistentFlags().StringVarP(&filepath, "filepath", "f", "", "<signing in user list file>.json")
	toCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 500, "-l <limit>, default 500")
	toCmd.MarkFlagRequired("filepath")

	toCmd.Example = "stress-test talent signin -f ~/Downloads/2W-user.json"

	rootCmd.AddCommand(toCmd)
}

var toCmd = &cobra.Command{
	Use:   "talent <action: signin>",
	Short: "Talent optimization test",
	Long:  `Talent optimization test`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if limit < 1 {
			log.Fatalf("limit <%v> must greater than 0\n", limit)
		}

		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			log.Fatalf("file not exits <%s>", filepath)
		}

		fileBuffer, err := ioutil.ReadFile(filepath)
		if err != nil {
			log.Fatalf("open file failed - %v\n", err)
		}

		var userList []map[string]string

		json.Unmarshal(fileBuffer, &userList)

		userLength := len(userList)
		if userLength == 0 {
			log.Fatal("no user loaded")
		}

		fmt.Printf("loaded %v users \n", userLength)

		httpClient := NewHttpClientWithoutRedirect(true)
		s := client.NewStressClientWithConcurrentNumber(1, userLength)

		rateLimiter := ratelimit.New(limit)
		var index uint32 = 0

		s.Header()
		s.Run(func() error {
			rateLimiter.Take()
			tmpIndex := atomic.AddUint32(&index, 1)

			talent := new(talent.TalentObject)
			signErr := talent.SignIn(userList[tmpIndex-1], httpClient)
			// signErr := talent.Status(httpClient)
			return signErr
		})
	},
}
