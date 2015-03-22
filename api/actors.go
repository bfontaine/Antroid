package api

// This file declares the interface `ActorInterface` and a type that implements
// it, `Actor`. An actor is here an external processus we start and which run
// in parallel of the game server. We can communicate with it using
// its `Send(string)` and `Read()` methods.
//
// Each actor is used to run an external command. We use them to run plugins
// like AIs and GUIs (see `api/plugins.go`). All actors have two booleans
// attributes, `readable` and `writable`. While all of them should be writable
// (i.e. we can send them messages) some of them don't need to be readable
// (e.g. we don't get messages from GUIs).
//
// The `Stage` structure represents a pool of `Actors`. It can contain any
// number of actors, can broadcast messages to them and read back all they
// messages.
//
// We implement actors with goroutines and communicate using two channels (one
// for the input and one for the output). `Stage`s use a WaitGroup to
// synchronize all actors.

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
	// A special message we send to actors to tell them to stop
	stop = "STOP"
)

// ActorInterface is used to represent generic Actors
type ActorInterface interface {

	// Start the actor. This should create a goroutine.
	Start(*sync.WaitGroup)
	// Send a message to the actor (blocking)
	Send(string)
	// Read the next message from the actor (blocking)
	Read() string
}

// An Actor represents a basic actor. The *Actor type (pointer on Actor)
// implements ActorInterface. We don't use the type directly, we "inherit" from
// it using anonymous structs in `api/ai.go` and `api/plugins.go`.
type Actor struct {
	// The command to run
	cmd *exec.Cmd

	// input channel, one need to write on it to send messages to the actor
	input chan string
	// output channel, which can be used to read from the actor
	output chan string

	// We can always communicate with actors, but `readable` means we'll be
	// able to read the external command's output (stdout) and `writable` menas
	// we'll able to write to it (stdin).
	readable, writable bool
}

// NewActor returns a pointer on a new Actor object. `cmd` is the command to
// start. It won't be run until `Start(wg)` is called.
func NewActor(cmd *exec.Cmd, readable, writable bool) *Actor {
	return &Actor{
		cmd: cmd,

		input:  make(chan string),
		output: make(chan string),

		readable: readable,
		writable: writable,
	}
}

// Start this actor. This takes a WaitGroup as an argument, which will be
// incremented when the actor starts, and decremented when it ends. It's used
// to ensure all actors ended when we terminate a `Stage`.
func (a *Actor) Start(wg *sync.WaitGroup) {
	// just use the unexported .start(wg) method in a goroutine
	go a.start(wg)
}

// Send a message to this actor (blocking)
func (a *Actor) Send(m string) { a.input <- m }

// Receive a message from this actor (blocking)
func (a *Actor) Read() string { return <-a.output }

// errLog takes an error and prints it on stderr along with the actor's command
func (a *Actor) errLog(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s", a.cmd.Path, err)
}

// This is the main function of an actor. It binds I/O and starts an infinite
// loop in which it reads an input, send it to the underlying command, wait for
// it to write a line on stdout and send back this line on its output channel.
//
// If an actor is not readable it won't write anything on its output channel.
// If it's not writable it *will* read on its input channel but won't send
// these messages to the command. It needs to read to get special messages like
// "stop". This method returns an error but we'll never be able to get it since
// it's started in the goroutine. This is why we log it using `.errLog` (see
// above).
func (a *Actor) start(wg *sync.WaitGroup) (err error) {
	var stdin io.WriteCloser
	var stdout io.ReadCloser

	// notify the wait group when we're done
	if wg != nil {
		defer wg.Done()
	}

	// create a pipe for STDIN
	if stdin, err = a.cmd.StdinPipe(); err != nil {
		a.errLog(err)
		return
	}

	// create a pipe for STDOUT
	if stdout, err = a.cmd.StdoutPipe(); err != nil {
		a.errLog(err)
		return
	}

	// redirect STDERR on our own one
	a.cmd.Stderr = os.Stderr

	// close STDIN if we're not writable
	if !a.writable {
		stdin.Close()
	}

	// close STDOUT if we're not readable
	if !a.readable {
		stdout.Close()
	}

	// start the underlying command
	if err = a.cmd.Start(); err != nil {
		a.errLog(err)
		return
	}

	// bufferize our STDOUT pipe to be able to use higher level reading
	// methods
	stdoutReader := bufio.NewReader(stdout)

	var buf []byte
	var msg string

	// main loop
	for {
		// 1. read on our input channel
		msg = <-a.input

		// if it's the special "stop" message, break the loop
		if msg == stop {
			break
		}

		// 2. if we're writable, send the input on the command's STDIN
		if a.writable {
			if _, err = io.WriteString(stdin, msg); err != nil {
				a.errLog(err)
				break
			}
		}

		// 3. if we're readable...
		if a.readable {
			// 3.1 read one line on the command's STDOUT...
			if buf, err = stdoutReader.ReadSlice('\n'); err != nil {
				a.errLog(err)
				break
			}

			// 3.2 ...and send it on our output channel
			a.output <- string(buf)
		}
	}

	// close STDIN before waiting for the command to end
	if a.writable {
		stdin.Close()
	}

	// wait for the command to end
	a.cmd.Wait()

	// close our input/output channels. This means we can't start an actor
	// twice because the second time its channels will be closed, but we don't
	// do that anyway.
	close(a.input)
	close(a.output)

	return
}

// A Stage is an extensible set of Actors
type Stage struct {
	actors []ActorInterface

	// This is the Wait Group we use to synchronize all actors at the end
	wg *sync.WaitGroup
}

// NewStage returns a new, empty, stage
func NewStage() *Stage {
	return &Stage{
		wg: &sync.WaitGroup{},
	}
}

// AddActor adds a new actor on the stage. All calls to this method should
// happen before the call to `.Start()`.
func (s *Stage) AddActor(a ActorInterface) {
	s.actors = append(s.actors, a)
}

// Start starts all actors
func (s *Stage) Start() {
	for _, a := range s.actors {
		s.wg.Add(1)
		a.Start(s.wg)
	}
}

// SendAll sends the string to all actors
func (s *Stage) SendAll(msg string) {
	for _, a := range s.actors {
		a.Send(msg)
	}
}

// ReadAll returns a slice of all actors' output. Trailing newlines are
// removed.
func (s *Stage) ReadAll() []string {
	msgs := make([]string, len(s.actors))

	for i, a := range s.actors {
		msgs[i] = strings.TrimSuffix(a.Read(), "\n")
	}

	return msgs
}

// Stop sends a special "stop" message to all actors and wait for them to
// terminate.
func (s *Stage) Stop() {
	s.SendAll(stop)
	s.wg.Wait()
}
