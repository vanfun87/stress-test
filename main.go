package main

import (
	"net/http"

	"github.com/ginkgoch/stress-test/pkg/client"
	"github.com/ginkgoch/stress-test/pkg/templates"
)

func main() {
	s := client.NewStressClientWithConcurrentNumber(20000, 40)

	s.Header()
	s.Run(func() error {
		request, _ := http.NewRequest("GET", "http://localhost:3000/version", nil)
		err := templates.HttpGet(request)
		return err
	})
}
