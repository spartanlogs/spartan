package codecs

import (
	"fmt"

	"github.com/lfkeitel/spartan/src/common"
)

type Codec interface {
	Format(*common.Event) []byte
}

var registeredCodecs map[string]Codec

func registerCodec(name string, c Codec) {
	if registeredCodecs == nil {
		registeredCodecs = make(map[string]Codec)
	}
	registeredCodecs[name] = c
}

func NewCodec(name string) (Codec, error) {
	c, exists := registeredCodecs[name]
	if !exists {
		return nil, fmt.Errorf("Codec %s doesn't exist", name)
	}
	return c, nil
}
