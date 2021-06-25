package stress

import (
	"fmt"
	"time"
)

type Stress struct {
	Number     int
	Concurrent int
	Limitation int
}

func (s *Stress) Run(taskFunc func() error) {
	printPrepareMessage(s)

	startTime := uint64(time.Now().UnixNano())
	endTime := uint64(startTime)

	successNum := 0
	failureNum := 0
	var maxTime, minTime, processTime uint64 = 0, 0, 0

	for i := 0; i < s.Number; i++ {
		r := runSingleTask(taskFunc)
		processTime += r.ProcessingTime

		if r.Success {
			successNum++
		} else {
			failureNum++
		}

		if maxTime == 0 {
			maxTime = r.ProcessingTime
		} else if r.ProcessingTime > maxTime {
			maxTime = r.ProcessingTime
		}

		if minTime == 0 {
			minTime = r.ProcessingTime
		} else if r.ProcessingTime < minTime {
			minTime = r.ProcessingTime
		}
	}
	endTime = uint64(time.Now().UnixNano())
	totalTime := endTime - startTime

	fmt.Printf("tasks takes %d ms\n", totalTime/uint64(time.Microsecond))
	fmt.Printf("process takes %d ms\n", processTime/uint64(time.Microsecond))
	fmt.Printf("max process takes %d ms\n", maxTime/uint64(time.Microsecond))
	fmt.Printf("min process takes %d ms\n", minTime/uint64(time.Microsecond))

}

func NewStress(number int, concurrent int, limitation int) *Stress {
	return &Stress{
		Number:     number,
		Concurrent: concurrent,
		Limitation: limitation,
	}
}

func NewStressWithNumber(number int) *Stress {
	return NewStress(number, 1, 0)
}

func NewStressWithConcurrentNumber(number int, concurrent int) *Stress {
	return NewStress(number, concurrent, 0)
}

func printPrepareMessage(s *Stress) {
	msg := fmt.Sprintf("%d task(s) ready to run with %d thread(s)", s.Number, s.Concurrent)

	if s.Limitation > 0 {
		msg += fmt.Sprintf("%s, with %d task(s) limitation per second", msg, s.Limitation)
	}

	fmt.Println(msg)
	fmt.Println()
}

func runSingleTask(taskFunc func() error) *TaskResult {
	startTime := time.Now()
	err := taskFunc()
	endTime := time.Now()

	processingTime := uint64(endTime.UnixNano() - startTime.UnixNano())
	return &TaskResult{
		Success:        err == nil,
		ProcessingTime: processingTime,
	}
}
