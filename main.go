package main

import (
	"time"

	"github.com/ginkgoch/stress-test/pkg/client"
)

func main() {
	s := client.NewStressClientWithNumber(1000)
	s.Header()
	s.RunSync(func() error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})
}
