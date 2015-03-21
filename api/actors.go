package api

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

const (
	stop = "STOP"
)

type ActorInterface interface {
	Start(*sync.WaitGroup)
	Send(string)
	Read() string
}

type Actor struct {
	cmd *exec.Cmd

	input  chan string
	output chan string

	readable, writable bool
}

func NewActor(cmd *exec.Cmd, readable, writable bool) *Actor {
	return &Actor{
		cmd: cmd,

		input:  make(chan string),
		output: make(chan string),

		readable: readable,
		writable: writable,
	}
}

func (a *Actor) Start(wg *sync.WaitGroup) {
	go a.start(wg)
}

func (a *Actor) Send(m string) {
	a.input <- m
}
func (a *Actor) Read() string { return <-a.output }

func (a *Actor) ErrLog(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s", a.cmd.Path, err)
}

func (a *Actor) start(wg *sync.WaitGroup) (err error) {
	var stdin io.WriteCloser
	var stdout io.ReadCloser

	// notify the wait group when we're done
	if wg != nil {
		defer wg.Done()
	}

	// create a pipe for STDIN
	if stdin, err = a.cmd.StdinPipe(); err != nil {
		a.ErrLog(err)
		return
	}

	// create a pipe for STDOUT
	if stdout, err = a.cmd.StdoutPipe(); err != nil {
		a.ErrLog(err)
		return
	}

	// redirect STDERR on our own one
	a.cmd.Stderr = os.Stderr

	if !a.writable {
		stdin.Close()
	}

	if !a.readable {
		stdout.Close()
	}

	// start the underlying command
	if err = a.cmd.Start(); err != nil {
		a.ErrLog(err)
		return
	}

	// bufferize our STDOUT pipe to be able to use higher level reading
	// methods
	stdoutReader := bufio.NewReader(stdout)

	var buf []byte
	var msg string

	for {
		fmt.Printf("%s\n", a.input)
		msg = <-a.input

		// special "stop" message
		if msg == stop {
			break
		}

		if a.writable {
			if _, err = io.WriteString(stdin, msg); err != nil {
				a.ErrLog(err)
				break
			}
		}

		if a.readable {
			if buf, err = stdoutReader.ReadSlice('\n'); err != nil {
				a.ErrLog(err)
				break
			}

			a.output <- string(buf)
		}
	}

	if a.writable {
		stdin.Close()
	}

	// wait for the command to stop
	a.cmd.Wait()

	// close our input/output channels
	close(a.input)
	close(a.output)

	return
}

type Stage struct {
	actors []ActorInterface

	wg *sync.WaitGroup
}

func NewStage() *Stage {
	return &Stage{
		wg: &sync.WaitGroup{},
	}
}

func (s *Stage) AddActor(a ActorInterface) {
	s.actors = append(s.actors, a)
}

func (s *Stage) Start() {
	for _, a := range s.actors {
		s.wg.Add(1)
		a.Start(s.wg)
	}
}

func (s *Stage) SendAll(msg string) {
	for _, a := range s.actors {
		a.Send(msg)
	}
}

func (s *Stage) ReadAll() []string {
	msgs := make([]string, len(s.actors))

	for i, a := range s.actors {
		msgs[i] = strings.TrimSuffix(a.Read(), "\n")
	}

	return msgs
}

func (s *Stage) Stop() {
	s.SendAll(stop)
	s.wg.Wait()
}
