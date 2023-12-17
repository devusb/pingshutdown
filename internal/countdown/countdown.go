package countdown

import (
	"time"
)

type Countdown struct {
	endTime       time.Time
	Duration      time.Duration
	runningStatus bool
	timer         *time.Timer
}

func (c *Countdown) Start() {
	if !c.runningStatus {
		c.runningStatus = true
		c.timer = time.NewTimer(c.Duration)
		c.endTime = time.Now().Add(c.Duration)
	}
}

func (c *Countdown) StartAfterFunc(f func()) {
	if !c.runningStatus {
		c.runningStatus = true
		c.timer = time.AfterFunc(c.Duration, f)
		c.endTime = time.Now().Add(c.Duration)
	}
}

func (c *Countdown) Stop() {
	if c.runningStatus {
		c.runningStatus = false
		c.timer.Stop()
		c.endTime = time.Now()
	}
}

func (c *Countdown) EndTime() time.Time {
	return c.endTime
}

func (c *Countdown) Status() bool {
	return c.runningStatus
}
