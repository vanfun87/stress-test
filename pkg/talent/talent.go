package talent

import (
	"fmt"
	"net/http"

	"github.com/ginkgoch/stress-test/pkg/templates"
)

const (
	serviceEndPoint = "https://talent.test.moblab-us.cn/api/1"
	// serviceEndPoint = "http://talent:3000/api/1/zhilian"
	signInUrl = "/zhilian/login"
	statusUrl = "/status"
)

type TalentObject struct {
	Cookie string
}

func (talent *TalentObject) Status(httpClient *http.Client) error {
	request, err := http.NewRequest("GET", formalizeUrl(statusUrl), nil)
	if err != nil {
		return err
	}

	err = templates.HttpGet(request, httpClient)
	return err
}

func (talent *TalentObject) SignIn(user map[string]string, httpClient *http.Client) error {
	request, err := http.NewRequest("GET", formalizeUrl(signInUrl), nil)
	if err != nil {
		return err
	}

	query := request.URL.Query()
	for key := range user {
		query.Add(key, user[key])
	}
	query.Add("accessId", "111111")

	request.URL.RawQuery = query.Encode()
	res, err := httpClient.Do(request)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if _, err = templates.ConsumeResponse(res); err != nil {
		return err
	}

	for _, cookie := range res.Cookies() {
		if cookie.Name == "this.sid" {
			cookie := cookie.String()
			cookie = fmt.Sprintf("%s;%s", cookie, "undefined")
			talent.Cookie = cookie
			break
		}
	}

	return nil
}

func formalizeUrl(url string) string {
	return fmt.Sprintf("%s%s", serviceEndPoint, url)
}
