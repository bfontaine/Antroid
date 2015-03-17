package api

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"
	"strings"
)

type AI struct {
	c      *exec.Cmd
	Input  chan string
	Output chan string
}

func NewAI(name string, arg ...string) *AI {
	return &AI{
		c:      exec.Command(name, arg...),
		Input:  make(chan string),
		Output: make(chan string),
	}
}

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

	for {
		if _, err = io.WriteString(stdin, <-ai.Input); err != nil {
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

type AIPool struct {
	ais []*AI
}

func NewAIPool() *AIPool {
	return &AIPool{
		ais: []*AI{},
	}
}

func (pool *AIPool) AddAI(name string, args ...string) {
	pool.ais = append(pool.ais, NewAI(name, args...))
}

func (pool *AIPool) Start() {
	for _, ai := range pool.ais {
		go ai.Start()
	}
}

func (pool *AIPool) SendMessage(msg string) {
	for _, ai := range pool.ais {
		ai.Input <- msg
	}
}

func (pool *AIPool) GetResponse() (resp string) {
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

	return buf.String()
}
