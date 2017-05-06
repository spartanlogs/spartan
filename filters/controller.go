package filters

import (
	"time"

	"github.com/spartanlogs/spartan/event"
	tomb "gopkg.in/tomb.v2"
)

type controllerOption func(*Controller)

type Controller struct {
	pipelines []*Pipeline
	batchSize int
	t         tomb.Tomb
	out       chan<- *event.Event
	flushers  []Flushable
}

func newController(options ...controllerOption) *Controller {
	c := &Controller{
		batchSize: 150,
		pipelines: make([]*Pipeline, 0, 1),
		flushers:  make([]Flushable, 0),
	}

	for _, o := range options {
		o(c)
	}
	return c
}

func withBatchSize(batch int) controllerOption {
	return func(c *Controller) { c.batchSize = batch }
}

func (c *Controller) Start(in <-chan *event.Event, out chan<- *event.Event) {
	c.out = out
	for _, pipeline := range c.pipelines {
		pipeline.Start(c.batchSize, in, out)
	}
	if len(c.flushers) > 0 {
		c.t.Go(c.flush)
	}
}

func (c *Controller) Close() error {
	// Kill and wait for the flusher to finish
	c.t.Kill(nil)
	c.t.Wait()

	// Close all pipelines
	for _, pipeline := range c.pipelines {
		pipeline.Close()
	}
	return nil
}

func (c *Controller) addFlusher(flusher Flushable) {
	c.flushers = append(c.flushers, flusher)
}

func (c *Controller) flush() error {
	timerDur := 5 * time.Second
	timer := time.NewTimer(timerDur) // Use timer incase filters take too long to flush
	for {
		select {
		case <-c.t.Dying():
			return nil
		case <-timer.C:
			for _, filter := range c.flushers {
				for _, e := range filter.Flush() {
					c.out <- e
				}
			}
			timer.Reset(timerDur)
		}
	}
}
