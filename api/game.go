package api

// GameID is a game id
type GameID string

// Game represents a game
type Game struct {
	Identifier GameID
	Spec       *GameSpec
}
