package api

// This file describes plugins, which are like AIs in `api/ais.go` except that
// they don't produce any output. A plugin is thus an external command that
// "listen" to our game server.
//
// See `api/actors.go` to understand how actors work

import (
	"os/exec"
)

// A Listener is just a non-readable actor
type Listener struct{ *Actor }

// NewListener returns a pointer on a new Listener struct
func NewListener(name string, arg ...string) *Listener {
	return &Listener{
		Actor: NewActor(exec.Command(name, arg...), false, true),
	}
}

// A ListenersPool is just a wrapper around a Stage that contains Listeners
// only.
type ListenersPool struct{ Stage }

// NewListenersPool returns a new, empty, ListenersPool
func NewListenersPool() *ListenersPool {
	return &ListenersPool{
		Stage: *NewStage(),
	}
}

// AddListener adds a new Listener to the pool
func (pool *ListenersPool) AddListener(name string, args ...string) {
	pool.AddActor(NewListener(name, args...))
}
