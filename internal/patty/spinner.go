package patty

import (
	"fmt"
	"os"
	"time"
)

type Spinner struct {
	stop    chan bool
	done    chan bool
	message string
}

func NewSpinner(message string) *Spinner {
	return &Spinner{
		stop:    make(chan bool),
		done:    make(chan bool),
		message: message,
	}
}

func (s *Spinner) Start() {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	go func() {
		i := 0
		for {
			select {
			case <-s.stop:
				fmt.Print("\r" + clearLine(s.message) + "\r")
				s.done <- true
				return
			default:
				frame := frames[i%len(frames)]
				fmt.Printf("\r%s %s", frame, s.message)
				os.Stdout.Sync()
				time.Sleep(100 * time.Millisecond)
				i++
			}
		}
	}()
}

func (s *Spinner) Stop() {
	s.stop <- true
	<-s.done
}

func clearLine(msg string) string {
	// clear the line by printing spaces
	spaces := ""
	for i := 0; i < len(msg)+10; i++ {
		spaces += " "
	}
	return spaces
}
