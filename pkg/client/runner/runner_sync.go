package runner

import (
	"sync"
	"time"
)

func RunSync(num int, ch chan<- *TaskResult, wg *sync.WaitGroup, taskFunc func() error) {
	for i := 0; i < num; i++ {
		r := runSingleTask(taskFunc)
		ch <- r
	}

	wg.Done()
}

func runSingleTask(taskFunc func() error) *TaskResult {
	startTime := time.Now()
	err := taskFunc()
	endTime := time.Now()

	processingTime := uint64(endTime.UnixNano() - startTime.UnixNano())
	return &TaskResult{
		Success:     err == nil,
		ProcessTime: processingTime,
	}
}
