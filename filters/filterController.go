package filters

import (
	"context"
	"fmt"
	"time"

	"github.com/lfkeitel/spartan/event"

	tomb "gopkg.in/tomb.v2"
)

// A FilterController is responsible for collecting a batch of events from inputs
// and start a chain of Filters to process the batch. Events are then sent to
// outputs.
type FilterController struct {
	start     FilterWrapper
	batchSize int
	t         tomb.Tomb
	in        <-chan *event.Event
	out       chan<- *event.Event
}

// NewFilterController creates a new controller using start as the root Filter
// and batchSize as the number of events to queue before processing.
func NewFilterController(start FilterWrapper, batchSize int) *FilterController {
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

		// Weird bugfix for Windows when shutting down in PowerShell.
		if currentBatch == 0 {
			if stopping {
				return nil
			}
			continue
		}

		// Can happen do to the go routine dying
		// or the batch timeout was exceeded.
		if currentBatch < f.batchSize {
			batch = batch[:currentBatch]
		}

		fmt.Printf("Processing batch of %d\n", len(batch))
		start := time.Now()
		batch = f.start.Run(context.TODO(), batch)
		fmt.Println(time.Since(start))

		for _, event := range batch {
			f.out <- event
		}

		if stopping {
			return nil
		}
	}
}
