package main

import (
	"time"

	stress "github.com/ginkgoch/go-stress/pkg/client"
)

func main() {
	s := stress.NewStressClientWithNumber(100)
	s.Header()
	r := stress.RunSerial(s.Number, func() error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})

	r.Print()
}
