package client

import (
	"fmt"
	"sync"
	"time"

	"github.com/ginkgoch/stress-test/pkg/client/runner"
)

type StressTestClient struct {
	Number     int
	Concurrent int
	Limitation int
}

func NewStressClient(number int, concurrent int, limitation int) *StressTestClient {
	return &StressTestClient{
		Number:     number,
		Concurrent: concurrent,
		Limitation: limitation,
	}
}

func NewStressClientWithNumber(number int) *StressTestClient {
	return NewStressClient(number, 1, 0)
}

func NewStressClientWithConcurrentNumber(number int, concurrent int) *StressTestClient {
	return NewStressClient(number, concurrent, 0)
}

func (s *StressTestClient) Header() {
	msg := fmt.Sprintf("%d task(s) ready to run with %d thread(s)", s.Number, s.Concurrent)

	if s.Limitation > 0 {
		msg += fmt.Sprintf("%s, with %d task(s) limitation per second", msg, s.Limitation)
	}

	fmt.Println(msg)
	fmt.Println()
}

func (s *StressTestClient) RunSync(taskFunc func() error) {
	startTime := uint64(time.Now().UnixNano())

	successNum := 0
	failureNum := 0
	var maxTime, minTime, processTime uint64 = 0, 0, 0

	ch := make(chan *runner.TaskResult, 1000)
	wg := new(sync.WaitGroup)

	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
			}
		}
	}()

	wg.Add(1)
	go runner.RunSync(s.Number, ch, wg, taskFunc)
	wg.Wait()
}
