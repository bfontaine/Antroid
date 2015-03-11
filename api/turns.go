package api

const (
	// LEFT is the command to turn to the left (don't move)
	LEFT = "left"
	// RIGHT is the command to turn to the right (don't move)
	RIGHT = "right"
	// FORWARD is the command to move forward
	FORWARD = "forward"
	// REST is the command to rest (don't move)
	REST = "rest"
)

// Command is a command we give to an ant
type Command struct {
	Ant int
	Cmd string
}

type Commands []Command

type Turn struct {
	Number int

	// TODO observations, etc
}
