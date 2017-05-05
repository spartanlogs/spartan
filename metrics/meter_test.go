package metrics

import (
	"sync"
	"testing"
)

func TestMeterMark(t *testing.T) {
	m := NewMeter(SMA)
	m.Mark(1)

	if c := m.Count(); c != 1 {
		t.Errorf("Wrong count. Expected %d, got %d", 1, c)
	}
}

func markMeter(m *Meter, count int) {
	for i := 0; i < 100; i++ {
		m.Mark(1)
	}
}

func TestMeterMarkGoroutines(t *testing.T) {
	m := NewMeter(SMA)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			markMeter(m, 1000)
			wg.Done()
		}()
	}
	wg.Wait()

	if c := m.Count(); c != 1000 {
		t.Errorf("Wrong count. Expected %d, got %d", 1000, c)
	}
}
