package metrics

import (
	"sync"
	"time"
)

type simpleMovingAverage struct {
	sync.Mutex
	values   []int64
	interval int64
	index    int
}

func newSMA(dur, interval time.Duration) *simpleMovingAverage {
	durSec := dur.Nanoseconds() / int64(time.Second)
	intervalSec := interval.Nanoseconds() / int64(time.Second)

	return &simpleMovingAverage{
		interval: intervalSec,
		values:   make([]int64, 1, durSec/intervalSec),
	}
}

func (a *simpleMovingAverage) Add(v int64) {
	a.Lock()
	a.values[a.index] = a.values[a.index] + v
	a.Unlock()
}

func (a *simpleMovingAverage) Clear() {
	a.Lock()
	a.values = make([]int64, 1, cap(a.values))
	a.index = 0
	a.Unlock()
}

func (a *simpleMovingAverage) Tick() {
	a.Lock()
	a.index = (a.index + 1) % cap(a.values)
	if len(a.values) < cap(a.values) {
		a.values = append(a.values, 0)
	}
	a.values[a.index] = 0
	a.Unlock()
}

func (a *simpleMovingAverage) Rate() float64 {
	var num, count int64

	a.Lock()
	for _, v := range a.values {
		num += v
		count++
	}
	a.Unlock()

	return float64(num) / float64(count) / float64(a.interval)
}
