package templates

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var client *http.Client

func init() {
	// tr := &http.Transport{
	// 	DialContext: (&net.Dialer{
	// 		Timeout:   30 * time.Second,
	// 		KeepAlive: 30 * time.Second,
	// 	}).DialContext,
	// 	MaxIdleConns:        0,                // 最大连接数,默认0无穷大
	// 	MaxIdleConnsPerHost: 1,                // 对每个host的最大连接数量(MaxIdleConnsPerHost<=MaxIdleConns)
	// 	IdleConnTimeout:     90 * time.Second, // 多长时间未使用自动关闭连接
	// 	TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	// }
	client = &http.Client{
		// Transport: tr,
	}
}

func HttpGet(request *http.Request) error {
	res, err := client.Do(request)
	if err != nil {
		return err
	}

	if res.StatusCode < 200 || res.StatusCode >= 400 {
		return fmt.Errorf("status code <%d> error", res.StatusCode)
	}

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return nil
}
