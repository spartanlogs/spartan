package outputs

import (
	"fmt"

	"github.com/lfkeitel/spartan/event"

	tomb "gopkg.in/tomb.v2"
)

// A OutputController is responsible for collecting a batch of events from filter
// pipelines and starting a chain of Outputs to process the batch. Events are
// considered process once they've been through the output chain.
type OutputController struct {
	start     Output
	batchSize int
	t         tomb.Tomb
	in        <-chan *event.Event
	out       chan<- *event.Event
}

// NewOutputController creates a new controller using start as the root Output
// and batchSize as the number of events to queue before processing.
func NewOutputController(start Output, batchSize int) *OutputController {
	return &OutputController{
		start:     start,
		batchSize: batchSize,
	}
}

// Start creates a go routine where the controller will start to wait for
// and collect events for processing. The in channel is used to collect Events
// from filter pipelines.
func (o *OutputController) Start(in chan *event.Event) error {
	o.in = in
	o.t.Go(o.run)
	return nil
}

// Close will gracefully shutdown the Controller. Collection from the input channel
// is immediately stopped and all in-flight events are processed and then the
// controller go routine exits.
func (o *OutputController) Close() error {
	o.t.Kill(nil)
	return o.t.Wait()
}

func (o *OutputController) run() error {
	fmt.Println("Output Pipeline started")
	for {
		select {
		case <-o.t.Dying():
			return nil
		default:
		}

		currentBatch := 0
		batch := make([]*event.Event, o.batchSize)
		stopping := false

	CURRENT:
		for currentBatch < o.batchSize {
			select {
			case event := <-o.in:
				batch[currentBatch] = event
				currentBatch++
			case <-o.t.Dying():
				stopping = true
				break CURRENT
			}
		}

		//fmt.Println("Processing batch")
		o.start.Run(batch)

		if stopping {
			return nil
		}
	}
}
