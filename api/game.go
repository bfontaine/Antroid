package api

type GameId string

type Game struct {
	Identifier GameId
	Spec       *GameSpec
}
