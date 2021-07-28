package lib

import (
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

//DoAction do action
func DoAction(count int, thread int, speed int, doFunc func(int) (interface{}, error)) (rsList []interface{}) {
	if speed == 0 {
		return doActionThread(count, thread, doFunc)
	}
	return doActionSpeed(count, speed, doFunc)
}

func task2(c chan interface{}, index int, doFunc func(int) (interface{}, error)) {
	v, err := doFunc(index)
	if err != nil {
		c <- nil
		logrus.Infoln("work task error:", err)
		return
	}
	c <- v
}

func doActionSpeed(count int, speed int, doFunc func(int) (interface{}, error)) (rsList []interface{}) {
	ticker := time.NewTicker(time.Duration(1000000000/speed) * time.Nanosecond)
	logrus.Infof("ticker speed %d\n", speed)
	c := make(chan interface{}, 100)
	defer close(c)
	index, endCount := 0, 0
	for {
		select {
		case <-ticker.C:
			if index >= count {
				ticker.Stop()
				logrus.Infoln("ticker.Stop")
			} else {
				logrus.Infoln("task start ", index)
				go task2(c, index, doFunc)
				index++
			}
		case r := <-c:
			if r != nil {
				rsList = append(rsList, r)
			}
			endCount++
			if endCount >= count {
				logrus.Infoln("all task done")
				return rsList
			}
		}
	}
}

func doActionThread(count int, threadNumber int, doFunc func(int) (interface{}, error)) (rsList []interface{}) {

	c := make(chan interface{}, 100)
	defer close(c)
	endCount := 0
	var index int32 = -1
	for i := 0; i < threadNumber; i++ {
		go func() {
			for current := atomic.AddInt32(&index, 1); current < int32(count); current = atomic.AddInt32(&index, 1) {
				task2(c, int(current), doFunc)
			}
		}()
	}
	for r := range c {
		if r != nil {
			rsList = append(rsList, r)
		}
		endCount++
		if endCount >= count {
			logrus.Infoln("all task done")
			return
		}
	}
	return
}
