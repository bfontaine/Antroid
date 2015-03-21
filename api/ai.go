package api

import (
	"os/exec"
	"strings"
)

// An AI struct represents an external AI command
type AI struct {
	*Actor
}

// NewAI returns a pointer on a new AI
func NewAI(name string, arg ...string) *AI {
	return &AI{
		Actor: NewActor(exec.Command(name, arg...), true, true),
	}
}

// An AIPool is a pool of multiple AIs
type AIPool struct {
	Stage
}

// NewAIPool returns a new empty AIPool
func NewAIPool() *AIPool {
	return &AIPool{
		Stage: *NewStage(),
	}
}

// AddAI adds another AI to the pool
func (pool *AIPool) AddAI(name string, args ...string) {
	pool.AddActor(NewAI(name, args...))
}

// GetCommandResponse reads the messages from all AIs
func (pool *AIPool) GetCommandResponse() (resp Commands) {
	return Commands(strings.Join(pool.ReadAll(), ","))
}
