package progress

import (
	"fmt"
	"io"
	"sync"
	"time"
)

type Spinner struct {
	writer io.Writer
	frames []string
	delay  time.Duration

	done chan struct{}
	wg   sync.WaitGroup
}

func New(writer io.Writer) *Spinner {
	return &Spinner{
		writer: writer,
		frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		delay:  80 * time.Millisecond,
	}
}

func (s *Spinner) Start(message string) {
	s.done = make(chan struct{})
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		frame := 0

		for {
			select {
			case <-s.done:
				return
			default:
				fmt.Fprintf(s.writer, "\r%s %s", message, s.frames[frame])
				frame = (frame + 1) % len(s.frames)
				time.Sleep(s.delay)
			}
		}
	}()
}

func (s *Spinner) Stop() {
	if s.done == nil {
		return
	}

	close(s.done)
	s.wg.Wait()

	fmt.Fprint(s.writer, "\r\033[K")
}
