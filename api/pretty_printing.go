package api

import (
	"fmt"
)

func (g Game) String() string {
	return fmt.Sprintf("Game %s, created on %s by %s (%s)",
		g.Identifier, g.CreationDate, g.Creator, g.Teaser)
}

func (g GameStatus) String() string {
	return fmt.Sprintf("Game %s, created on %s by %s (%s), turn %d (%s)",
		g.Identifier, g.CreationDate, g.Creator, g.Teaser, g.Turn, g.Status)
}

func (cmd Command) String() string {
	return string(cmd)
}

func (cmds Commands) String() string {
	return string(cmds)
}

func (resp baseResponse) String() string {
	return fmt.Sprintf(`{"status": "%s", "response": %s}`,
		resp.Status, resp.Response)
}

func (t Turn) String() string {
	return fmt.Sprintf("%d", t.Number)
}
