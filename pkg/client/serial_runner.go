package client

import (
	"time"
)

func RunSync(num int, taskFunc func() error) *SerialTaskResult {
	startTime := uint64(time.Now().UnixNano())
	endTime := uint64(startTime)

	successNum := 0
	failureNum := 0
	var maxTime, minTime, processTime uint64 = 0, 0, 0

	for i := 0; i < num; i++ {
		r := runSingleTask(taskFunc)
		processTime += r.ProcessTime

		if r.Success {
			successNum++
		} else {
			failureNum++
		}

		if maxTime == 0 {
			maxTime = r.ProcessTime
		} else if r.ProcessTime > maxTime {
			maxTime = r.ProcessTime
		}

		if minTime == 0 {
			minTime = r.ProcessTime
		} else if r.ProcessTime < minTime {
			minTime = r.ProcessTime
		}
	}
	endTime = uint64(time.Now().UnixNano())
	totalTime := endTime - startTime

	serialResult := &SerialTaskResult{
		SuccessNum:  successNum,
		FailureNum:  failureNum,
		ProcessTime: processTime,
		SerialTime:  totalTime,
		MaxTime:     maxTime,
		MinTime:     minTime,
	}

	return serialResult
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
