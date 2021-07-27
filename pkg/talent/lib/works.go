package lib

import "time"

type DelayWork struct {
	delay time.Time
	work  func()
}

var delayWorkPools []chan *DelayWork

func InitWorkPools(number int) {
	for i := 0; i < number; i++ {
		delayWorkPools = append(delayWorkPools, make(chan *DelayWork, 20000))
	}
}

func SendWork(work func(), delay int) {
	if delay <= 0 {
		work()
	}
	if delay > len(delayWorkPools) {
		delay = len(delayWorkPools)
	}
	delayWorkPools[delay-1] <- &DelayWork{delay: time.Now().Add(time.Duration(delay) * time.Second), work: work}
}

func RunDelayWorkPool() {
	for _, works := range delayWorkPools {
		go func() {
			for work := range works {
				if t := work.delay.Sub(time.Now()); t > 0 {
					time.Sleep(t)
				}
				work.work()
			}
		}()
	}
}

var SendWorkPool = make(chan func(), 200)

func RunSendWorkPool() {
	for i := 0; i < 2; i++ {
		go func() {
			for work := range SendWorkPool {
				work()
			}
		}()
	}
}

//func createCases(chs ...chan int) []reflect.SelectCase {
//}
