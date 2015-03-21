package api

import (
	"os/exec"
)

type Listener struct {
	*Actor
}

func NewListener(name string, arg ...string) *Listener {
	return &Listener{
		Actor: NewActor(exec.Command(name, arg...), false, true),
	}
}

// An AIPool is a pool of multiple AIs
type ListenersPool struct {
	Stage
}

// NewAIPool returns a new empty AIPool
func NewListenersPool() *ListenersPool {
	return &ListenersPool{
		Stage: *NewStage(),
	}
}

// AddAI adds another AI to the pool
func (pool *ListenersPool) AddListener(name string, args ...string) {
	pool.AddActor(NewListener(name, args...))
}
