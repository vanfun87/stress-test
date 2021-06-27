package statistics

import (
	"fmt"
	"sync"
	"time"

	"github.com/ginkgoch/stress-test/pkg/client/runner"
)

type ResultStatistics struct {
	ConcurrentNum int
	StartTime,
	SuccessNum,
	FailureNum,
	MaxTime,
	MinTime,
	RunningTime,
	ProcessTime uint64
}

func NewResultStatistics(concurrentNum int) *ResultStatistics {
	if concurrentNum == 0 {
		concurrentNum = 1
	}

	return &ResultStatistics{
		ConcurrentNum: concurrentNum,
		StartTime:     0,
		SuccessNum:    0,
		FailureNum:    0,
		MaxTime:       0,
		MinTime:       0,
		RunningTime:   0,
		ProcessTime:   0,
	}
}

func (s *ResultStatistics) Watch(ch <-chan *runner.TaskResult, wg *sync.WaitGroup) {
	defer wg.Done()

	s.StartTime = uint64(time.Now().UnixNano())

	stopCh := make(chan bool)
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				go s.PrintTableRow()
			case <-stopCh:
				return
			}
		}
	}()

	s.PrintTableHeader()
	for r := range ch {
		s.Append(r)
	}

	stopCh <- true
	s.PrintTableRow()
}

func (s *ResultStatistics) Append(r *runner.TaskResult) {
	s.ProcessTime += r.ProcessTime
	s.RunningTime = uint64(time.Now().UnixNano()) - s.StartTime

	if r.Success {
		s.SuccessNum++
	} else {
		s.FailureNum++
	}

	if s.MaxTime == 0 || r.ProcessTime > s.MaxTime {
		s.MaxTime = r.ProcessTime
	}

	if s.MinTime == 0 || r.ProcessTime < s.MinTime {
		s.MinTime = r.ProcessTime
	}
}

func (s *ResultStatistics) PrintTableHeader() {
	fmt.Println("─────────┬─────────┬─────────┬──────────┬──────────┬──────────┬──────────")
	fmt.Println(" 耗时(s) │  成功数 │  失败数 │     qps  │ 最长耗时 │ 最短耗时 │ 平均耗时 ")
	fmt.Println("─────────┼─────────┼─────────┼──────────┼──────────┼──────────┼──────────")
}

func (s *ResultStatistics) PrintTableRow() {
	processTime := s.ProcessTime
	if processTime == 0 {
		processTime = 1
	}
	row := fmt.Sprintf(" %7d │ %7d │ %7d │ %8.2f │ %8.2f │ %8.2f │ %8.2f ",
		s.RunningTime/1e9,
		s.SuccessNum,
		s.FailureNum,
		float64(s.SuccessNum*uint64(s.ConcurrentNum)*1e9)/float64(processTime),
		float64(s.MaxTime)/1e6,
		float64(s.MinTime)/1e6,
		float64(s.ProcessTime)/1e6/float64(s.SuccessNum+s.FailureNum),
	)
	fmt.Println(row)
}
