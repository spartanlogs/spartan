package codecs

import (
	"github.com/lfkeitel/spartan/common"
)

type DotCodec struct{}

func init() {
	register("dot", &DotCodec{})
}

func (c *DotCodec) Format(e *common.Event) []byte {
	return []byte{'.'}
}
