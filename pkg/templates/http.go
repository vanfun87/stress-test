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

func SendRequest(request *http.Request, client *http.Client, handler func(*http.Response)) ([]byte, error) {
	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if handler != nil {
		handler(res)
	}

	if res.StatusCode < 200 || res.StatusCode >= 400 {
		return nil, fmt.Errorf("status code <%d> error", res.StatusCode)
	}

	buffer, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}
