package codecs

import (
	"io"
	"testing"

	"github.com/spartanlogs/spartan/event"
)

type testCodec struct{}

func newTestCodec() (Codec, error)                                     { return &testCodec{}, nil }
func (c *testCodec) Encode(e *event.Event) []byte                      { return nil }
func (c *testCodec) EncodeWriter(io.Writer, <-chan *event.Event)       {}
func (c *testCodec) Decode(data []byte) (*event.Event, error)          { return nil, nil }
func (c *testCodec) DecodeReader(r io.Reader, out chan<- *event.Event) {}

func TestRegistration(t *testing.T) {
	defer func() {
		delete(registeredCodecInits, "test")
		if r := recover(); r == nil {
			t.Fatal("Register didn't panic from dual registration")
		}
	}()

	Register("test", newTestCodec)

	if len(registeredCodecInits) != 1 {
		t.Fatalf("registeredCodecInits length. Expected %d, got %d", 1, len(registeredCodecInits))
	}

	// Should panic
	Register("test", newTestCodec)
}

func TestNewCodec(t *testing.T) {
	Register("test", newTestCodec)

	c, err := New("test")
	if err != nil {
		t.Fatalf("Got error creating codec: %s", err)
	}

	if _, ok := c.(*testCodec); !ok {
		t.Fatalf("New codec not correct. Got %T", c)
	}
}
