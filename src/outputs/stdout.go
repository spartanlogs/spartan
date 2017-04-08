package outputs

import (
	"fmt"

	"github.com/lfkeitel/spartan/src/codecs"
	"github.com/lfkeitel/spartan/src/common"
)

func NewStdoutOutput(next Output) Output {
	c, _ := codecs.NewCodec("json")
	return func(batch []*common.Event) []*common.Event {
		for _, event := range batch {
			if event != nil {
				fmt.Printf("%s\n", c.Format(event))
			}
		}
		return next(batch)
	}
}
