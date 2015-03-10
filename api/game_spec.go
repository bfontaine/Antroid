package api

import (
	"strings"
)

// GameSpec represents all the parameters needed to define a game
type GameSpec struct {
	Public        bool
	Players       []string
	Description   string
	Pace          int
	Turns         int
	AntsPerPlayer int
	MaxPlayers    int
	MinPlayers    int
	InitialEnergy int
	InitialAcid   int
}

func (gs *GameSpec) toParams() GameSpecParams {
	gsp := GameSpecParams{
		Teaser:          gs.Description,
		Pace:            gs.Pace,
		NbTurn:          gs.Turns,
		NbAntPerPlayer:  gs.AntsPerPlayer,
		NbPlayer:        gs.MaxPlayers,
		MinimalNbPlayer: gs.MinPlayers,
		InitialEnergy:   gs.InitialEnergy,
		InitialAcid:     gs.InitialAcid,
	}

	if gs.Public {
		gsp.Users = "all"
	} else {
		gsp.Users = strings.Join(gs.Players, ",")
	}

	return gsp
}

type intRange struct {
	min int
	max int
}

// IntRange creates a new intRange
func IntRange(min, max int) intRange {
	return intRange{min: min, max: max}
}

// Include checks if an int is included in the range
func (r intRange) Include(n int) bool {
	return r.min <= n && n <= r.max
}

// constants for the v0 API
var (
	paceRange          = IntRange(1, 100)
	turnsRange         = IntRange(1, 100000)
	antsPerPlayerRange = IntRange(1, 42)
	playersRange       = IntRange(1, 42)
	initialEnergyRange = IntRange(1, 1000)
	initialAcidRange   = IntRange(1, 1000)
)

// Validate checks that the spec validates against the spec spec
func (gs *GameSpec) Validate() bool {
	nbUsers := len(gs.Players)

	// games are either public or private. In the later case they must have 1+
	// players.
	if (gs.Public && nbUsers > 0) || (!gs.Public && nbUsers == 0) {
		return false
	}

	// TODO check if the API accepts empty teasers

	if !paceRange.Include(gs.Pace) ||
		!turnsRange.Include(gs.Turns) ||
		!antsPerPlayerRange.Include(gs.AntsPerPlayer) ||
		!playersRange.Include(gs.MinPlayers) ||
		!playersRange.Include(gs.MaxPlayers) ||
		gs.MinPlayers > gs.MaxPlayers ||
		!initialEnergyRange.Include(gs.InitialEnergy) ||
		!initialAcidRange.Include(gs.InitialAcid) {
		return false
	}

	return true
}
