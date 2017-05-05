package metrics

import (
	"sync/atomic"
	"time"
)

type Meter struct {
	count int64
	rates []averager
	stop  chan struct{}
}

func NewMeter(avgType AverageType) *Meter {
	m := &Meter{
		count: 0,
		rates: make([]averager, 3),
		stop:  make(chan struct{}),
	}

	switch avgType {
	case SMA:
		m.rates[0] = newSMA(1*time.Minute, 5*time.Second)
		m.rates[1] = newSMA(5*time.Minute, 5*time.Second)
		m.rates[2] = newSMA(15*time.Minute, 5*time.Second)
	}

	go m.tick(5 * time.Second)

	return m
}

func (m *Meter) tick(rate time.Duration) {
	ticker := time.NewTicker(rate)
	for {
		select {
		case <-ticker.C:
			m.rates[0].Tick()
			m.rates[1].Tick()
			m.rates[2].Tick()
		case <-m.stop:
			ticker.Stop()
			return
		}
	}
}

func (m *Meter) Clear() {
	atomic.StoreInt64(&m.count, 0)
	m.rates[0].Clear()
	m.rates[1].Clear()
	m.rates[2].Clear()
}

func (m *Meter) Mark(val int64) {
	atomic.AddInt64(&m.count, 1)
	m.rates[0].Add(val)
	m.rates[1].Add(val)
	m.rates[2].Add(val)
}

func (m *Meter) Count() int64 {
	return atomic.LoadInt64(&m.count)
}

func (m *Meter) Stop() {
	m.stop <- struct{}{}
}

func (m *Meter) OneMinuteRate() float64 {
	return m.rates[0].Rate()
}

func (m *Meter) FiveMinuteRate() float64 {
	return m.rates[1].Rate()
}

func (m *Meter) FifteenMinuteRate() float64 {
	return m.rates[2].Rate()
}
