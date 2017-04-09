package outputs

import (
	"fmt"

	tomb "gopkg.in/tomb.v2"

	"github.com/lfkeitel/spartan/src/common"
)

type Output interface {
	SetNext(Output)
	Run([]*common.Event)
}

type OutputController struct {
	start     Output
	batchSize int
	t         tomb.Tomb
	in        <-chan *common.Event
	out       chan<- *common.Event
}

func NewOutputController(start Output, batchSize int) *OutputController {
	return &OutputController{
		start:     start,
		batchSize: batchSize,
	}
}

func (o *OutputController) Start(in chan *common.Event) error {
	o.in = in
	o.t.Go(o.run)
	return nil
}

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
		batch := make([]*common.Event, o.batchSize)
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

		fmt.Println("Processing batch")
		o.start.Run(batch)

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
