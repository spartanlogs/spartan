package inputs

import (
	"fmt"
	"time"

	"github.com/hpcloud/tail"
	"github.com/lfkeitel/spartan/event"
	"gopkg.in/tomb.v2"
)

type FileInput struct {
	path string
	t    tomb.Tomb
	out  chan<- *event.Event
}

func NewFileInput(file string) *FileInput {
	return &FileInput{
		path: file,
	}
}

func (i *FileInput) Start(out chan<- *event.Event) error {
	i.out = out
	i.t.Go(i.run)
	return nil
}

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

		t, err := tail.TailFile(i.path, tail.Config{
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
