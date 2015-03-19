package api

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

const (
	stop = "STOP"
)

// An AI struct represents an external AI command
type AI struct {
	c      *exec.Cmd
	Input  chan string
	Output chan string
}

// NewAI returns a pointer on a new AI
func NewAI(name string, arg ...string) *AI {
	return &AI{
		c:      exec.Command(name, arg...),
		Input:  make(chan string),
		Output: make(chan string),
	}
}

// Start starts the AI. You should start it in a goroutine.
func (ai *AI) Start(wg *sync.WaitGroup) (err error) {

	if wg != nil {
		defer wg.Done()
	}

	stdin, err := ai.c.StdinPipe()

	// Note: we might want to redirect each AI's stderr in a specific file
	ai.c.Stderr = os.Stderr

	if err != nil {
		return
	}

	stdout, err := ai.c.StdoutPipe()

	if err != nil {
		return
	}

	defer stdin.Close()

	if err = ai.c.Start(); err != nil {
		return
	}

	stdoutReader := bufio.NewReader(stdout)

	var buf []byte
	var msg string

	for {

		msg = <-ai.Input

		if msg == stop {
			break
		}

		if _, err = io.WriteString(stdin, msg); err != nil {
			break
		}

		if buf, err = stdoutReader.ReadSlice('\n'); err != nil {
			break
		}

		ai.Output <- string(buf)
	}

	ai.c.Wait()

	close(ai.Input)
	close(ai.Output)

	return
}

// An AIPool is a pool of multiple AIs
type AIPool struct {
	ais []*AI

	wg *sync.WaitGroup
}

// NewAIPool returns a new empty AIPool
func NewAIPool() *AIPool {
	return &AIPool{
		ais: []*AI{},
		wg:  &sync.WaitGroup{},
	}
}

// AddAI adds another AI to the pool
func (pool *AIPool) AddAI(name string, args ...string) {
	pool.ais = append(pool.ais, NewAI(name, args...))
}

// Start starts all AIs in separate goroutines
func (pool *AIPool) Start() {
	for _, ai := range pool.ais {
		pool.wg.Add(1)
		go ai.Start(pool.wg)
	}
}

// SendMessage sends a message to all AIs
func (pool *AIPool) SendMessage(msg string) {
	for _, ai := range pool.ais {
		ai.Input <- msg
	}
}

// GetCommandResponse reads the messages from all AIs
func (pool *AIPool) GetCommandResponse() (resp Commands) {
	var buf bytes.Buffer

	first := true

	for _, ai := range pool.ais {
		if !first {
			buf.WriteString(",")
		}
		line := <-ai.Output
		buf.WriteString(strings.TrimSuffix(line, "\n"))
		first = false
	}

	return Commands(buf.String())
}

// Stop stops all AIs in the pool and wait for them to terminate
func (pool *AIPool) Stop() {
	pool.SendMessage(stop)
	pool.wg.Wait()
}
