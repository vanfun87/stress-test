package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/ginkgoch/stress-test/pkg/client"
	"github.com/ginkgoch/stress-test/pkg/client/runner"
	"github.com/ginkgoch/stress-test/pkg/talent"
	"github.com/spf13/cobra"
	"go.uber.org/ratelimit"
)

var (
	filepath       string
	limit          int
	serverEndpoint string
	debug          bool
	stage          int
	useQps         bool
	game           bool
)

func init() {
	toCmd.PersistentFlags().StringVarP(&filepath, "filepath", "f", "", "<signing in user list file>.json")
	toCmd.PersistentFlags().StringVarP(&serverEndpoint, "serverEndpoint", "u", talent.DefaultServiceEndPoint, "https://talent.test.moblab-us.cn/api/1")
	toCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 500, "-l <limit>, default 500")
	toCmd.PersistentFlags().IntVarP(&stage, "stage", "t", 0, "-t <stage>, default 0")
	toCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "-d, default false")
	toCmd.PersistentFlags().BoolVarP(&useQps, "qps", "q", false, "-q, default false")
	toCmd.PersistentFlags().BoolVarP(&game, "game", "g", false, "-g, default false")
	toCmd.MarkFlagRequired("filepath")

	toCmd.Example = "stress-test talent -f ~/Downloads/2W-user.json"

	rootCmd.AddCommand(toCmd)
}

var toCmd = &cobra.Command{
	Use:   "talent",
	Short: "Talent optimization test",
	Long:  `Talent optimization test`,
	Args:  cobra.NoArgs,
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

		var httpClient *http.Client
		if ParseBool(keepAlive) {
			httpClient = NewHttpClientWithoutRedirect(true)
		}

		if debug {
			debugErr := executeSingleTask(userList[0], httpClient, nil)
			if debugErr != nil {
				log.Fatal(debugErr)
			}
		} else {
			executeStressTest(userList, httpClient)
		}
	},
}

func executeStressTest(userList []map[string]string, httpClient *http.Client) {
	s := client.NewStressClientWithConcurrentNumber(1, len(userList))

	rateLimiter := ratelimit.New(limit)
	var index uint32 = 0

	s.Header()

	if !useQps {
		s.RunSingleTaskWithRateLimiter(rateLimiter, func() error {
			tmpIndex := atomic.AddUint32(&index, 1)

			user := userList[tmpIndex-1]

			debugErr := executeSingleTask(user, httpClient, nil)
			return debugErr
		})
	} else {
		s.RunMultiTasksWithRateLimiter(rateLimiter, func(ch chan<- *runner.TaskResult) error {
			tmpIndex := atomic.AddUint32(&index, 1)

			user := userList[tmpIndex-1]

			debugErr := executeSingleTask(user, httpClient, ch)
			return debugErr
		})
	}
}

func executeSingleTask(user map[string]string, httpClient *http.Client, ch chan<- *runner.TaskResult) (err error) {
	if httpClient == nil {
		httpClient = NewHttpClientWithoutRedirect(false)
	}

	talentObj := talent.NewTalentObject(serverEndpoint)

	if stage == -1 {
		t1 := time.Now()
		err = talentObj.Status(httpClient)
		enqueueMetrics(&t1, err, ch)

		if err != nil {
			return err
		} else if debug {
			fmt.Println("debug - status success")
		}

		return
	}

	i := 0
	if stage == 0 || stage > i {
		t1 := time.Now()
		err = talentObj.SignIn(user, httpClient)
		enqueueMetrics(&t1, err, ch)

		if err != nil {
			return err
		} else if debug {
			fmt.Println("debug - sign in success with cookie", talentObj.Cookie)
		}

		i++
	}

	if stage == 0 || stage > i {
		t1 := time.Now()
		err = talentObj.Information(httpClient)
		enqueueMetrics(&t1, err, ch)

		if err != nil {
			return err
		} else if debug {
			fmt.Println("debug - information success", talentObj.UserId)
		}

		i++
	}

	if stage == 0 || stage > i {
		t1 := time.Now()
		err = talentObj.StartGame("competitive_math", httpClient)
		enqueueMetrics(&t1, err, ch)

		if err != nil {
			return err
		} else if debug {
			fmt.Println("debug - start game success", talentObj.GameConfig)
		}

		i++
	}

	if game && (stage == 0 || stage > i) {
		err = talentObj.PlayGame("competitive_math")
		if err != nil {
			return err
		} else if debug {
			fmt.Println("debug - play game success")
		}

		i++
	}

	if stage == 0 || stage > i {
		t1 := time.Now()
		err = talentObj.StopGame("competitive_math", httpClient)
		enqueueMetrics(&t1, err, ch)

		if err != nil {
			return err
		} else if debug {
			fmt.Println("debug - stop game success")
		}

		i++
	}

	return
}

func enqueueMetrics(startTime *time.Time, err error, ch chan<- *runner.TaskResult) {
	d := time.Since(*startTime).Nanoseconds()
	if ch != nil {
		ch <- &runner.TaskResult{
			Success:     err == nil,
			ProcessTime: uint64(d),
		}
	}
}
