package metrics

import (
	"math"
	"sync"
	"sync/atomic"
	"time"
)

var (
	m1Alpha  = 1.0 - math.Exp(-5.0/60/1)
	m5Alpha  = 1.0 - math.Exp(-5.0/60/5)
	m15Alpha = 1.0 - math.Exp(-5.0/60/15)
)

type ewma struct {
	sync.Mutex
	interval    int64
	rate        float64
	uncounted   int64
	alpha       float64
	initialized bool
}

func newEWMA(alpha float64, interval time.Duration) *ewma {
	return &ewma{
		alpha:    alpha,
		interval: interval.Nanoseconds() / int64(time.Second),
	}
}

func (e *ewma) Add(v int64) {
	atomic.AddInt64(&e.uncounted, v)
}
func (e *ewma) Tick() {
	e.Lock()
	instantRate := float64(e.uncounted) / float64(e.interval)
	e.uncounted = 0

	if !e.initialized {
		e.rate = float64(instantRate)
		e.initialized = true
	} else {
		e.rate = float64(e.rate) + e.alpha*(instantRate-float64(e.rate))
	}
	e.Unlock()
}
func (e *ewma) Rate() float64 {
	var rate float64
	e.Lock()
	rate = e.rate
	e.Unlock()
	return rate
}
func (e *ewma) Clear() {
	e.Lock()
	e.initialized = false
	e.rate = 0.0
	e.uncounted = 0
	e.Unlock()
}
