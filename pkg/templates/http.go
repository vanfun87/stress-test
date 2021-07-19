package templates

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func HttpGet(request *http.Request, client *http.Client) error {
	res, err := client.Do(request)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 400 {
		return fmt.Errorf("status code <%d> error", res.StatusCode)
	}

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return nil
}
