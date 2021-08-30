package statistics

import (
	"sync"
	"time"

	"github.com/ginkgoch/stress-test/pkg/client/runner"
)

type TimeWindow struct {
	SizeInSec int
	Buckets   []*TimeBucket
	Locker    sync.RWMutex
}

func NewTimeWindow(windowSizeInSec int) *TimeWindow {
	return &TimeWindow{SizeInSec: windowSizeInSec}
}

type TimeBucket struct {
	Key                  int64
	StartTimeInNanoSec   int64
	EndTimeInNanoSec     int64
	ProcessTimeInNanoSec int64
	SuccessNum           int
	FailureNum           int
}

func (window *TimeWindow) Append(r *runner.TaskResult) {
	currentTime := time.Now().UnixNano()
	currentTimeInSec := currentTime / 1e9

	window.Locker.Lock()
	defer window.Locker.Unlock()

	// clean timeout buckets
	window.CleanTimeoutBuckets(currentTimeInSec - int64(window.SizeInSec))

	var bucket *TimeBucket
	if len(window.Buckets) > 0 {
		lastBucket := window.Buckets[len(window.Buckets)-1]
		if lastBucket.Key == currentTimeInSec {
			bucket = lastBucket
		}
	}

	if bucket == nil {
		bucket = &TimeBucket{Key: currentTimeInSec, StartTimeInNanoSec: currentTime}
		window.Buckets = append(window.Buckets, bucket)
	}

	bucket.EndTimeInNanoSec = currentTime

	if r.Success {
		bucket.SuccessNum++
	} else {
		bucket.FailureNum++
	}

	bucket.ProcessTimeInNanoSec += int64(r.ProcessTime)
}

func (window *TimeWindow) Info() (qps float64, speed float64) {
	window.Locker.RLock()
	defer window.Locker.RUnlock()

	var (
		startTime   int64
		endTime     int64
		duration    int64
		successSum  int64
		failureSum  int64
		processTime int64
	)

	for i := 0; i < len(window.Buckets); i++ {
		bucket := window.Buckets[i]
		if startTime == 0 || startTime > bucket.StartTimeInNanoSec {
			startTime = bucket.StartTimeInNanoSec
		}

		if endTime == 0 || endTime < bucket.EndTimeInNanoSec {
			endTime = bucket.EndTimeInNanoSec
		}

		successSum += int64(bucket.SuccessNum)
		failureSum += int64(bucket.FailureNum)
		processTime += bucket.ProcessTimeInNanoSec
	}

	duration = endTime - startTime
	totalSum := successSum + failureSum
	qps = float64(totalSum) / float64(duration/1e9)
	speed = float64(processTime/1e6) / float64(totalSum)
	return
}

func (window *TimeWindow) CleanTimeoutBuckets(clearTimeSinceInSec int64) {
	for len(window.Buckets) > 0 && window.Buckets[0].Key < clearTimeSinceInSec {
		window.Buckets = window.Buckets[1:]
	}
}
