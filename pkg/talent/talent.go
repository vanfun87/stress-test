package talent

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ginkgoch/stress-test/pkg/templates"
)

const (
	serviceEndPoint = "https://talent.test.moblab-us.cn/api/1/zhilian"
	signInUrl       = "/login"
)

type TalentObject struct {
	Cookie string
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

	cookie := strings.Split(res.Cookies()[0].String(), ";")[0]
	cookie = fmt.Sprintf("%s;%s", cookie, "undefined")
	talent.Cookie = cookie

	return nil
}

func formalizeUrl(url string) string {
	return fmt.Sprintf("%s%s", serviceEndPoint, url)
}
