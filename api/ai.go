package api

// This files describe how we handle external AIs. See `api/actors.go` for the
// low-level operations. An `AI` is just a wrapper for an `Actor` and an
// `AIPool` is a `Stage` that contains only `AI`s.

import (
	"os/exec"
	"strings"
)

// An AI is just an Actor
type AI struct {
	*Actor
}

// NewAI returns a pointer on a new AI, which is an actor that is both readable
// and writable.
func NewAI(name string, arg ...string) *AI {
	return &AI{
		Actor: NewActor(exec.Command(name, arg...), true, true),
	}
}

// An AIPool is a pool of multiple AIs. This is just a wrapper around a `Stage`
// which contains only `AI`s.
type AIPool struct {
	Stage
}

// NewAIPool returns a pointer on a new, empty, AIPool
func NewAIPool() *AIPool {
	return &AIPool{
		Stage: *NewStage(),
	}
}

// AddAI adds another AI to the pool
func (pool *AIPool) AddAI(name string, args ...string) {
	pool.AddActor(NewAI(name, args...))
}

// GetCommandResponse reads the messages from all AIs and return them all as a
// Commands object that can be sent to the remote server.
func (pool *AIPool) GetCommandResponse() (resp Commands) {
	return Commands(strings.Join(pool.ReadAll(), ","))
}
