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
		Teaser:            gs.Description,
		Pace:              gs.Pace,
		Nb_turn:           gs.Turns,
		Nb_ant_per_player: gs.AntsPerPlayer,
		Nb_player:         gs.MaxPlayers,
		Minimal_nb_player: gs.MinPlayers,
		Initial_energy:    gs.InitialEnergy,
		Initial_acid:      gs.InitialAcid,
	}

	if gs.Public {
		gsp.Users = "+"
	} else {
		gsp.Users = strings.Join(gs.Players, ",")
	}

	return gsp
}

// constants for the v0 API
const (
	minPace          = 1
	maxPace          = 100
	minTurns         = 1
	maxTurns         = 100000
	minAntsPerPlayer = 1
	maxAntsPerPlayer = 42
	minPlayers       = 1
	maxPlayers       = 42
	minInitialEnergy = 1
	maxInitialEnergy = 1000
	minInitialAcid   = 1
	maxInitialAcid   = 1000
)

// Validate checks that the spec validates against the spec spec
func (gs *GameSpec) Validate() bool {
	nbUsers := len(gs.Players)

	// games are either public or private. In the later case they must have 1+
	// players.
	if (gs.Public && nbUsers > 0) || (!gs.Public && nbUsers == 0) {
		return false
	}

	if gs.Pace < minPace ||
		gs.Pace > maxPace ||
		gs.Turns < minTurns ||
		gs.Turns > maxTurns ||
		gs.AntsPerPlayer < minAntsPerPlayer ||
		gs.AntsPerPlayer > maxAntsPerPlayer ||
		gs.MaxPlayers < minPlayers ||
		gs.MaxPlayers > maxPlayers ||
		gs.MinPlayers < minPlayers ||
		gs.MinPlayers > gs.MaxPlayers ||
		gs.InitialEnergy < minInitialEnergy ||
		gs.InitialEnergy > maxInitialEnergy ||
		gs.InitialAcid < minInitialAcid ||
		gs.InitialAcid > maxInitialAcid {
		return false
	}

	return true
}
