package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
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
	filepath string
	limit    int
	debug    bool
	stage    int
	useQps   bool
	game     bool
)

func init() {
	toCmd.PersistentFlags().StringVarP(&filepath, "filepath", "f", "", "<signing in user list file>.json")
	toCmd.PersistentFlags().StringVarP(&talent.ServiceEndpoint, "serverEndpoint", "u", talent.DefaultServiceEndpoint, "https://talent.test.moblab-us.cn/api/1")
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

	tmpIndex := atomic.AddUint32(&index, 1)
	user := userList[tmpIndex-1]

	if !useQps {
		s.RunSingleTaskWithRateLimiter("talent", rateLimiter, func() error {
			debugErr := executeSingleTask(user, httpClient, nil)
			return debugErr
		})
	} else {
		s.RunMultiTasksWithRateLimiter("talent", rateLimiter, func(ch chan<- *runner.TaskResult) error {
			debugErr := executeSingleTask(user, httpClient, ch)
			return debugErr
		})
	}
}

func executeSingleStep(i int, action string, talentObj *talent.TalentObject, ch chan<- *runner.TaskResult, handler func() error) (int, error) {
	if stage == 0 || stage > i {
		t1 := time.Now()
		err := handler()
		enqueueMetrics(action, &t1, err, ch)

		if err != nil {
			return i, err
		} else if debug {
			fmt.Printf("debug - %s success: %v\n", action, talentObj.String())
		}

		i++
	}

	return i, nil
}

func executeSingleTask(user map[string]string, httpClient *http.Client, ch chan<- *runner.TaskResult) (err error) {
	if httpClient == nil {
		httpClient = NewHttpClientWithoutRedirect(false)
	}

	talentObj := talent.NewTalentObject()

	if stage == -1 {
		t1 := time.Now()
		err = talentObj.Status(httpClient)
		enqueueMetrics("status", &t1, err, ch)

		if err != nil {
			return err
		} else if debug {
			fmt.Printf("debug - status success")
		}

		return
	}

	i := 0
	if i, err = executeSingleStep(i, "sign-in", talentObj, ch, func() error {
		return talentObj.SignIn(user, httpClient)
	}); err != nil {
		return
	}

	if i, err = executeSingleStep(i, "information", talentObj, ch, func() error {
		return talentObj.Information(httpClient)
	}); err != nil {
		return
	}

	if i, err = executeSingleStep(i, "start-game", talentObj, ch, func() error {
		rand.Seed(time.Now().UnixNano())
		sleeping := rand.Intn(4000) + 1000
		time.Sleep(time.Duration(sleeping) * time.Millisecond)
		return talentObj.StartGame("competitive_math", httpClient)
	}); err != nil {
		return
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

	if _, err = executeSingleStep(i, "stop-game", talentObj, ch, func() error {
		sleeping := rand.Intn(1000) + 500
		time.Sleep(time.Duration(sleeping) * time.Millisecond)
		return talentObj.StopGame("competitive_math", httpClient)
	}); err != nil {
		return
	}

	return
}

func enqueueMetrics(name string, startTime *time.Time, err error, ch chan<- *runner.TaskResult) {
	endTime := time.Now()
	d := endTime.Sub(*startTime).Nanoseconds()

	if ch != nil {
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}

		ch <- &runner.TaskResult{
			Success:     err == nil,
			ProcessTime: uint64(d),
			StartTime:   uint64(startTime.UnixNano()),
			EndTime:     uint64(endTime.UnixNano()),
			Category:    name,
			Err:         errMsg,
		}
	}
}
