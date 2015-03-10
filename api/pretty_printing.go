package api

import "fmt"

func (g Game) String() string {
	return fmt.Sprintf("Game %s, created on %s by %s (%s)",
		g.Identifier, g.CreationDate, g.Creator, g.Teaser)
}

func (g GameStatus) String() string {
	return fmt.Sprintf("Game %s, created on %s by %s (%s), turn %d (%s)",
		g.Identifier, g.CreationDate, g.Creator, g.Teaser, g.Turn, g.Status)
}
