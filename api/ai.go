package api

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"
	"strings"
)

const (
	STOP = "STOP"
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
func (ai *AI) Start() (err error) {

	stdin, err := ai.c.StdinPipe()

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

		if msg == STOP {
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
}

// NewAIPool returns a new empty AIPool
func NewAIPool() *AIPool {
	return &AIPool{
		ais: []*AI{},
	}
}

// AddAI adds another AI to the pool
func (pool *AIPool) AddAI(name string, args ...string) {
	pool.ais = append(pool.ais, NewAI(name, args...))
}

// Start starts all AIs in separate goroutines
func (pool *AIPool) Start() {
	for _, ai := range pool.ais {
		go ai.Start()
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

func (pool *AIPool) Stop() {
	pool.SendMessage(STOP)
}
