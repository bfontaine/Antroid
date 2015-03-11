package api

import (
	"bytes"
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
	return fmt.Sprintf("%d:%s", cmd.Ant, cmd.Cmd)
}

func (cmds Commands) String() string {
	var buf bytes.Buffer

	for i, cmd := range cmds {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(cmd.String())
	}

	return buf.String()
}
