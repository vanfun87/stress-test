package main

import (
	"time"

	"github.com/ginkgoch/stress-test/pkg/client"
)

func main() {
	s := client.NewStressClientWithConcurrentNumber(1000, 4)
	s.Header()
	s.Run(func() error {
		time.Sleep(20 * time.Millisecond)
		return nil
	})
}
