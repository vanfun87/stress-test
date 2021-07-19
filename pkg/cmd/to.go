package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/ginkgoch/stress-test/pkg/templates"
	"github.com/spf13/cobra"
)

const (
	serviceEndPoint = "https://talent.test.moblab-us.cn/api/1/zhilian"
	signInUrl       = "/login"
)

var (
	filepath string
)

func init() {
	toCmd.PersistentFlags().StringVarP(&filepath, "filepath", "f", "", "<signing in user list file>.json")

	toCmd.MarkFlagRequired("filepath")

	toCmd.Example = "stress-test to signin -f ~/Downloads/2W-user.json"

	rootCmd.AddCommand(toCmd)
}

var toCmd = &cobra.Command{
	Use:   "to <test>",
	Short: "Talent optimization test",
	Long:  `Talent optimization test`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
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

		httpClient := NewHttpClient(true)
		doSignIn(userList[0], httpClient)
	},
}

func doSignIn(user map[string]string, httpClient *http.Client) {
	request, _ := http.NewRequest("GET", formalizeUrl(signInUrl), nil)

	query := request.URL.Query()
	for key := range user {
		query.Add(key, user[key])
	}
	query.Add("accessId", "111111")

	request.URL.RawQuery = query.Encode()

	fmt.Println(request.URL.String())
	_, err := templates.SendRequest(request, httpClient, func(res *http.Response) {
		// cookie := res.Header[http.CanonicalHeaderKey("set-cookie")]
		fmt.Println("cookie is", res.Header)
	})

	if err != nil {
		fmt.Println("request failed", err)
	} else {
		fmt.Println("sign in successfully")
	}
}

func formalizeUrl(url string) string {
	return fmt.Sprintf("%s%s", serviceEndPoint, url)
}
