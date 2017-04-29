package inputs

import (
	"errors"
	"fmt"
	"time"

	"github.com/hpcloud/tail"
	"github.com/lfkeitel/spartan/event"
	"github.com/lfkeitel/spartan/utils"
	"gopkg.in/tomb.v2"
)

func init() {
	Register("file", newFileInput)
}

type fileConfig struct {
	path string
}

// A FileInput will read a file and optionally tail it. Each line is considered
// a separate event.
type FileInput struct {
	config *fileConfig
	t      tomb.Tomb
	out    chan<- *event.Event
}

func newFileInput(options utils.InterfaceMap) (Input, error) {
	i := &FileInput{
		config: &fileConfig{},
	}
	return i, i.setConfig(options)
}

func (i *FileInput) setConfig(options utils.InterfaceMap) error {
	if s, exists := options.GetOK("path"); exists {
		i.config.path = s.(string)
	} else {
		return errors.New("Path option required")
	}

	return nil
}

// Start the FileInput reading/tailing.
func (i *FileInput) Start(out chan<- *event.Event) error {
	i.out = out
	i.t.Go(i.run)
	return nil
}

// Close the FileInput
func (i *FileInput) Close() error {
	i.t.Kill(nil)
	return i.t.Wait()
}

func (i *FileInput) run() error {
	for {
		select {
		case <-i.t.Dying():
			return nil
		default:
		}

		t, err := tail.TailFile(i.config.path, tail.Config{
			Follow: true,
			ReOpen: true,
		})

		if err != nil {
			fmt.Println(err.Error())
			time.Sleep(500 * time.Millisecond)
			continue
		}

		var line *tail.Line

		for {
			select {
			case line = <-t.Lines:
				if line.Err != nil {
					fmt.Println(line.Err.Error())
					continue
				}
				i.out <- event.New(line.Text)
			case <-i.t.Dying():
				t.Stop()
				t.Cleanup()
				return nil
			}
		}
	}
}
