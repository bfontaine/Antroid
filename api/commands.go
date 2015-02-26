package api

var (
	LEFT    = "left"
	RIGHT   = "right"
	FORWARD = "forward"
	REST    = "rest"
)

type Command struct {
	Ant int
	Cmd string
}
