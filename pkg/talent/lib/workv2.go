package lib

import "time"

type DelayWork struct {
	delay time.Time
	work  func()
}

type DelayWorkPools struct {
	DelayWorkChan []chan *DelayWork
}

func (p *DelayWorkPools) InitWorkPools(number int) {
	for i := 0; i < number; i++ {
		p.DelayWorkChan = append(p.DelayWorkChan, make(chan *DelayWork, 100))
	}
}

func (p *DelayWorkPools) SendWork(work func(), delay int) {
	if delay <= 0 {
		work()
	}
	if delay > len(p.DelayWorkChan) {
		delay = len(p.DelayWorkChan)
	}
	p.DelayWorkChan[delay-1] <- &DelayWork{delay: time.Now().Add(time.Duration(delay) * time.Second), work: work}
}

func (p *DelayWorkPools) RunDelayWorkPool() {
	for _, works := range p.DelayWorkChan {
		go func(ch chan *DelayWork) {
			for work := range ch {
				if t := time.Until(work.delay); t > 0 {
					time.Sleep(t)
				}
				work.work()
			}
		}(works)
	}
}
