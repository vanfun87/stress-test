package main

import (
	"time"

	"github.com/ginkgoch/stress-test/pkg/client"
)

func main() {
	s := client.NewStressClientWithNumber(100)
	s.Header()
	r := s.RunSync(func() error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})

	r.Print()
}
