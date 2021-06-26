package main

import (
	"time"

	client "github.com/ginkgoch/stress-test/pkg/client"
)

func main() {
	s := client.NewStressClientWithNumber(100)
	s.Header()
	r := client.RunSync(s.Number, func() error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})

	r.Print()
}
