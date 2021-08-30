package lib

import (
	"fmt"
	"runtime"
	"time"

	"github.com/ginkgoch/stress-test/pkg/log"
)

type StopWatch struct {
	timeMap   map[string]time.Time
	user      string
	printFile bool
}

func NewStopWatch(userid string) StopWatch {
	return StopWatch{timeMap: map[string]time.Time{}, user: userid, printFile: false}
}
func getFileAndLine(depth int) (file string) {
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		file = "???"
		line = 0
	}
	return fmt.Sprintf("%s:%d:", file, line)
}
func (sw *StopWatch) Start(name string, tag string) {
	t := time.Now()
	sw.timeMap[name] = t
	filepath := ""
	if sw.printFile {
		filepath = getFileAndLine(2)
	}

	log.Printf("%s %-10s %-13s %s %-20s %s \n", filepath, sw.user, "start_time:", time.Now().Format("2006-01-02 15:04:05.000000"), name, tag)
}

func (sw *StopWatch) Get(name string, tag string) {
	t, ok := sw.timeMap[name]
	if !ok {
		return
	}

	filepath := ""
	if sw.printFile {
		filepath = getFileAndLine(2)
	}
	log.Printf("%s %-10s %-13s %10dms %s -> %s \n", filepath, sw.user, "past_time:", time.Since(t)/time.Millisecond, name, tag)
}

func (sw *StopWatch) End(name string, tag string) {
	t, ok := sw.timeMap[name]
	if !ok {
		return
	}
	filepath := ""
	if sw.printFile {
		filepath = getFileAndLine(2)
	}
	log.Printf("%s %-10s %-13s %10dms %s -> %s \n", filepath, sw.user, "end_past_time:", time.Since(t)/time.Millisecond, name, tag)
	delete(sw.timeMap, name)
}

func (sw *StopWatch) GetPastTime(name string) time.Duration {
	t, ok := sw.timeMap[name]
	if !ok {
		return 0
	}
	return time.Since(t)
}

func (sw *StopWatch) Log(name string, msg string) {
	filepath := ""
	if sw.printFile {
		filepath = getFileAndLine(2)
	}
	log.Printf("%s %-10s Log: %-10s %s \n", filepath, sw.user, name, msg)

}
