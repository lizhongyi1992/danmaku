package syncer

import "time"

type TimerTask struct {
	stopchan  chan struct{}
	perSecond int
	f         func()
	running   bool
	ticker    *time.Ticker
}

func NewTimerTask(perSecond int, f func()) *TimerTask {
	if perSecond < 1 {
		// limit
		perSecond = 1
	}
	p := &TimerTask{
		stopchan:  make(chan struct{}),
		perSecond: perSecond,
		f:         f,
		running:   false,
		ticker:    time.NewTicker(time.Duration(perSecond) * time.Second),
	}
	return p
}

func (p *TimerTask) Start() {
	p.running = true
	for p.running == true {
		select {
		case <-p.ticker.C:
			p.f()
		}
	}
}

func (p *TimerTask) Stop() {
	close(p.stopchan)
	p.running = false
}

func (p *TimerTask) StopChan() chan struct{} {
	return p.stopchan
}
