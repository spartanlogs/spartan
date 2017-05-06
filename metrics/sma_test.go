package metrics

import (
	"testing"
	"time"
)

func TestSMAAddRate(t *testing.T) {
	sma := newSMA(1*time.Minute, 5*time.Second)
	sma.Add(1000)
	sma.Tick()
	if rate := sma.Rate(); rate != 100.0 {
		t.Errorf("Wrong rate. Expected %f, got %f", 100.0, rate)
	}
}
