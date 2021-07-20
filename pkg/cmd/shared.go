package cmd

import (
	"net/http"
	"time"
)

var trueFlags []string = []string{"true", "t", "1"}

func ContainsStr(arr []string, target string) bool {
	for _, i := range arr {
		if i == target {
			return true
		}
	}

	return false
}

func ParseBool(v string) bool {
	if ContainsStr(trueFlags, v) {
		return true
	} else {
		return false
	}
}

func NewHttpClient(keepAlive bool) *http.Client {
	var tr *http.Transport

	if keepAlive {
		tr = &http.Transport{
			MaxIdleConnsPerHost: 1024,
			TLSHandshakeTimeout: 0 * time.Second,
		}
	} else {
		tr = new(http.Transport)
	}

	httpClient := &http.Client{Transport: tr}
	return httpClient
}

func NewHttpClientWithoutRedirect(keepAlive bool) *http.Client {
	var tr *http.Transport

	if keepAlive {
		tr = &http.Transport{
			MaxIdleConnsPerHost: 1024,
			TLSHandshakeTimeout: 0 * time.Second,
		}
	} else {
		tr = new(http.Transport)
	}

	httpClient := &http.Client{Transport: tr, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}}

	return httpClient
}
