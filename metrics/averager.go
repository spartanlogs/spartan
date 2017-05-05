package metrics

type averager interface {
	Add(v int64)
	Tick()
	Rate() float64
	Clear()
}

type AverageType int

const (
	SMA  AverageType = 1
	EWMA AverageType = 2
)
