package filters

import (
	"fmt"
	"time"

	"github.com/lfkeitel/spartan/config"
	"github.com/lfkeitel/spartan/event"
	"github.com/lfkeitel/spartan/metrics"
	"github.com/lfkeitel/spartan/utils"
)

func init() {
	Register("metrics", newMetricsFilter)
}

type metricsConfig struct {
	meter         string
	flushInterval time.Duration
}

var metricsConfigSchema = []config.Setting{
	{
		Name:    "meter",
		Type:    config.String,
		Default: "events",
	},
	{
		Name:    "flush_interval",
		Type:    config.Int,
		Default: 5,
	},
}

// A MetricsFilter is used to perform several different actions on an Event.
// See the documentation for the Mutate filter for more information.
type MetricsFilter struct {
	config    *metricsConfig
	lastFlush time.Time
	meter     *metrics.Meter
}

func newMetricsFilter(options utils.InterfaceMap) (Filter, error) {
	options = checkOptionsMap(options)
	f := &MetricsFilter{
		config:    &metricsConfig{},
		lastFlush: time.Now(),
		meter:     metrics.NewMeter(metrics.EWMA),
	}
	if err := f.setConfig(options); err != nil {
		return nil, err
	}
	return f, nil
}

func (f *MetricsFilter) setConfig(options utils.InterfaceMap) error {
	if err := config.VerifySettings(options, metricsConfigSchema); err != nil {
		return err
	}

	f.config.meter = options.Get("meter").(string)
	f.config.flushInterval = time.Duration(options.Get("flush_interval").(int)) * time.Second

	return nil
}

// Filter processes a batch.
func (f *MetricsFilter) Filter(batch []*event.Event, matchedFunc MatchFunc) []*event.Event {
	for range batch {
		f.meter.Mark(1)
	}
	return batch
}

func (f *MetricsFilter) Flush() []*event.Event {
	f.lastFlush.Add(5 * time.Second)

	if !f.needFlush() {
		return nil
	}

	e := event.New("")
	e.SetType("metrics")
	e.SetField(fmt.Sprintf("%s.count", f.config.meter), f.meter.Count())
	e.SetField(fmt.Sprintf("%s.rate_1m", f.config.meter), f.meter.OneMinuteRate())
	e.SetField(fmt.Sprintf("%s.rate_5m", f.config.meter), f.meter.FiveMinuteRate())
	e.SetField(fmt.Sprintf("%s.rate_15m", f.config.meter), f.meter.FifteenMinuteRate())
	return []*event.Event{e}
}

func (f *MetricsFilter) needFlush() bool {
	return time.Since(f.lastFlush).Seconds() > f.config.flushInterval.Seconds()
}
