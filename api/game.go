package api

import "fmt"

// GameID is a game id
type GameID string

// Game represents a game
type Game struct {
	Identifier   GameID
	CreationDate string `json:"creation_date"`
	Creator      string
	Teaser       string

	Spec *GameSpec
}

func (g Game) String() string {
	return fmt.Sprintf("Game %s, created on %s by %s (%s)",
		g.Identifier, g.CreationDate, g.Creator, g.Teaser)
}
