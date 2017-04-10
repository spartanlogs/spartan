package filters

import (
	"fmt"

	"github.com/lfkeitel/spartan/event"

	tomb "gopkg.in/tomb.v2"
)

type FilterController struct {
	start     Filter
	batchSize int
	t         tomb.Tomb
	in        <-chan *event.Event
	out       chan<- *event.Event
}

func NewFilterController(start Filter, batchSize int) *FilterController {
	return &FilterController{
		start:     start,
		batchSize: batchSize,
	}
}

func (f *FilterController) Start(in, out chan *event.Event) error {
	f.in = in
	f.out = out
	f.t.Go(f.run)
	return nil
}

func (f *FilterController) Close() error {
	f.t.Kill(nil)
	return f.t.Wait()
}

func (f *FilterController) run() error {
	fmt.Println("Filter Pipeline started")
	for {
		select {
		case <-f.t.Dying():
			return nil
		default:
		}

		currentBatch := 0
		batch := make([]*event.Event, f.batchSize)
		stopping := false

	CURRENT:
		for currentBatch < f.batchSize {
			select {
			case event := <-f.in:
				batch[currentBatch] = event
				currentBatch++
			case <-f.t.Dying():
				stopping = true
				break CURRENT
			}
		}

		fmt.Println("Processing batch")
		batch = f.start.Run(batch)

		for _, event := range batch {
			f.out <- event
		}

		if stopping {
			return nil
		}
	}
}

func checkOptionsMap(o map[string]interface{}) map[string]interface{} {
	if o == nil {
		o = make(map[string]interface{})
	}
	return o
}
