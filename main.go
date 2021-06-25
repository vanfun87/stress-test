package main

import (
	"time"

	"github.com/ginkgoch/go-stress/pkg/stress"
)

func main() {
	s := stress.NewStressWithNumber(4)
	s.Run(func() error {
		// for i := 0; i < 10000; i++ {

		// }

		time.Sleep(200 * time.Microsecond)
		return nil
	})
}
