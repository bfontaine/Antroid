package api

// This file adds `.String()` methods on our structs to get human-readable
// descriptions.

import "fmt"

func (g Game) String() string {
	return fmt.Sprintf("Game %s, created on %s by %s (%s)",
		g.Identifier, g.CreationDate, g.Creator, g.Teaser)
}

func (g GameStatus) String() string {
	return fmt.Sprintf("Game %s, created on %s by %s (%s), turn %d (%s)",
		g.Identifier, g.CreationDate, g.Creator, g.Teaser, g.Turn, g.Status)
}

func (cmds Commands) String() string {
	return string(cmds)
}

func (resp baseResponse) String() string {
	return fmt.Sprintf(`{"status": "%s", "response": %s}`,
		resp.Status, resp.Response)
}

func (t Turn) String() string {
	return fmt.Sprintf("turn %d", t.Number)
}

func (p Position) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

func (d Direction) String() string {
	return fmt.Sprintf("(%d, %d)", d.X, d.Y)
}
