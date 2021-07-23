package runner

import (
	"sync"
	"time"
)

func RunSync(name string, num int, ch chan<- *TaskResult, wg *sync.WaitGroup, taskFunc func() error) {
	defer wg.Done()

	for i := 0; i < num; i++ {
		r := runSingleTask(name, taskFunc)
		ch <- r
	}
}

func runSingleTask(name string, taskFunc func() error) *TaskResult {
	startTime := time.Now()
	err := taskFunc()
	endTime := time.Now()

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	processingTime := uint64(endTime.Sub(startTime).Nanoseconds())
	return &TaskResult{
		Success:     err == nil,
		ProcessTime: processingTime,
		StartTime:   uint64(startTime.UnixNano()),
		EndTime:     uint64(endTime.UnixNano()),
		Category:    name,
		Err:         errMsg,
	}
}

func RunSyncWithMultiTasks(num int, ch chan<- *TaskResult, wg *sync.WaitGroup, taskFunc func(ch chan<- *TaskResult) error) {
	defer wg.Done()

	for i := 0; i < num; i++ {
		taskFunc(ch)
	}
}
