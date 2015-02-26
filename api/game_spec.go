package api

import (
	"strings"
)

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
		teaser:            gs.Description,
		pace:              gs.Pace,
		nb_turn:           gs.Turns,
		nb_ant_per_player: gs.AntsPerPlayer,
		nb_player:         gs.MaxPlayers,
		minimal_nb_player: gs.MinPlayers,
		initial_energy:    gs.InitialEnergy,
		initial_acid:      gs.InitialAcid,
	}

	if gs.Public {
		gsp.users = "+"
	} else {
		gsp.users = strings.Join(gs.Players, ",")
	}

	return gsp
}

// constants for the v0 API
var (
	MIN_PACE            = 1
	MAX_PACE            = 100
	MIN_TURNS           = 1
	MAX_TURNS           = 100000
	MIN_ANTS_PER_PLAYER = 1
	MAX_ANTS_PER_PLAYER = 42
	MIN_PLAYERS         = 1
	MAX_PLAYERS         = 42
	MIN_INITIAL_ENERGY  = 1
	MAX_INITIAL_ENERGY  = 1000
	MIN_INITIAL_ACID    = 1
	MAX_INITIAL_ACID    = 1000
)

// Check that the spec validates against the spec spec
func (g *GameSpec) Validate() bool {
	nbUsers := len(g.Players)

	// games are either public or private. In the later case they must have 1+
	// players.
	if (g.Public && nbUsers > 0) || (!g.Public && nbUsers == 0) {
		return false
	}

	if g.Pace < MIN_PACE ||
		g.Pace > MAX_PACE ||
		g.Turns < MIN_TURNS ||
		g.Turns > MAX_TURNS ||
		g.AntsPerPlayer < MIN_ANTS_PER_PLAYER ||
		g.AntsPerPlayer > MAX_ANTS_PER_PLAYER ||
		g.MaxPlayers < MIN_PLAYERS ||
		g.MaxPlayers > MAX_PLAYERS ||
		g.MinPlayers < MIN_PLAYERS ||
		g.MinPlayers > g.MaxPlayers ||
		g.InitialEnergy < MIN_INITIAL_ENERGY ||
		g.InitialEnergy > MAX_INITIAL_ENERGY ||
		g.InitialAcid < MIN_INITIAL_ACID ||
		g.InitialAcid > MAX_INITIAL_ACID {
		return false
	}

	return true
}
