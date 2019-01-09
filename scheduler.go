package main

import (
	"time"
)

const INTERVAL_PERIOD time.Duration = 24 * time.Hour

type Scheduler struct {
	t *time.Timer
}

func getNextTickDuration() time.Duration {
	now := time.Now()
	nextTick := time.Date(now.Year(), now.Month(), now.Day(), config.JobTimeHour, config.JobTimeMinute, config.JobTimeSecond, 0, time.Local)
	for nextTick.Before(now) {
		nextTick = nextTick.Add(INTERVAL_PERIOD)
	}
	return nextTick.Sub(time.Now())
}

func (sc Scheduler) updateJobTicker() {
	sc.t.Reset(getNextTickDuration())
}

func newScheduler(task func()) {
	sc := Scheduler{time.NewTimer(getNextTickDuration())}
	for {
		<-sc.t.C
		task()
		sc.updateJobTicker()
	}
}
