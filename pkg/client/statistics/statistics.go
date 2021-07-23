package statistics

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ginkgoch/stress-test/pkg/client/runner"
	log "github.com/sirupsen/logrus"
)

const logFilepath = "./stress-test.log"

func initLogger(enableLogger bool) {
	if !enableLogger {
		return
	}

	log.SetFormatter(&log.JSONFormatter{})

	var file, err = os.OpenFile(logFilepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Could Not Open Log File : " + err.Error())
	}

	log.SetOutput(file)

	log.SetLevel(log.InfoLevel)
}

type ResultStatistics struct {
	ConcurrentNum int
	StartTime,
	SuccessNum,
	FailureNum,
	MaxTime,
	MinTime,
	RunningTime,
	ProcessTime uint64
	EnableLogger bool
}

func NewResultStatistics(concurrentNum int, enableLogger bool) *ResultStatistics {
	initLogger(enableLogger)

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
		EnableLogger:  enableLogger,
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
		if s.EnableLogger {
			log.Info(r)
		}
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
		// qps can also be more precise when no rate limiter involved
		// float64(s.SuccessNum*uint64(s.ConcurrentNum)*1e9) / float64(processTime)
		float64(s.SuccessNum*1e9)/float64(s.RunningTime),
		float64(s.MaxTime)/1e6,
		float64(s.MinTime)/1e6,
		float64(s.ProcessTime)/1e6/float64(s.SuccessNum+s.FailureNum),
	)
	fmt.Println(row)
}
