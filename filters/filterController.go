package filters

import (
	"fmt"

	"github.com/lfkeitel/spartan/event"

	tomb "gopkg.in/tomb.v2"
)

// A FilterController is responsible for collecting a batch of events from inputs
// and start a chain of Filters to process the batch. Events are then sent to
// outputs.
type FilterController struct {
	start     Filter
	batchSize int
	t         tomb.Tomb
	in        <-chan *event.Event
	out       chan<- *event.Event
}

// NewFilterController creates a new controller using start as the root Filter
// and batchSize as the number of events to queue before processing.
func NewFilterController(start Filter, batchSize int) *FilterController {
	return &FilterController{
		start:     start,
		batchSize: batchSize,
	}
}

// Start creates a go routine where the controller will start to wait for
// and collect events for processing. The in channel is used to collect Events
// from inputs. The out channel is where Events are sent to the outputs.
func (f *FilterController) Start(in, out chan *event.Event) error {
	f.in = in
	f.out = out
	f.t.Go(f.run)
	return nil
}

// Close will gracefully shutdown the Controller. Collection from the input channel
// is immediately stopped and all in-flight events are processed, sent to outputs, and
// then the controller go routine exits.
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

// checkOptionsMap ensures an option map is never nil.
func checkOptionsMap(o map[string]interface{}) map[string]interface{} {
	if o == nil {
		o = make(map[string]interface{})
	}
	return o
}
