package main

import (
	"time"

	stress "github.com/ginkgoch/go-stress/pkg/client"
)

func main() {
	s := stress.NewStressClientWithNumber(4)
	s.Run(func() error {
		// for i := 0; i < 10000; i++ {

		// }

		time.Sleep(200 * time.Microsecond)
		return nil
	})
}
