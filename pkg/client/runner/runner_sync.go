package runner

import (
	"sync"
	"time"
)

func RunSync(concurrentId int, num int, ch chan<- *TaskResult, wg *sync.WaitGroup, taskFunc func(i int) error) {
	defer wg.Done()

	for i := 0; i < num; i++ {
		r := runSingleTask(concurrentId, taskFunc)
		ch <- r
	}
}

func runSingleTask(concurrentId int, taskFunc func(i int) error) *TaskResult {
	startTime := time.Now()
	err := taskFunc(concurrentId)
	endTime := time.Now()

	processingTime := uint64(endTime.UnixNano() - startTime.UnixNano())
	return &TaskResult{
		Success:     err == nil,
		ProcessTime: processingTime,
	}
}
