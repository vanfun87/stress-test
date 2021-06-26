package client

import (
	"fmt"
	"sync"
	"time"

	"github.com/ginkgoch/stress-test/pkg/client/runner"
	"github.com/ginkgoch/stress-test/pkg/client/statistics"
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
	ch := make(chan *runner.TaskResult, 1000)
	wg := new(sync.WaitGroup)
	wgStatistics := new(sync.WaitGroup)

	st := statistics.NewResultStatistics()

	wgStatistics.Add(1)
	go st.Watch(ch, wgStatistics)

	wg.Add(1)
	go runner.RunSync(s.Number, ch, wg, taskFunc)
	wg.Wait()

	time.Sleep(1 * time.Millisecond)
	close(ch)

	wgStatistics.Wait()
}
