package api

import (
	"encoding/json"
)

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

// A GameStatus represents the current status of a game
type GameStatus struct {
	Game

	Score  map[string]int
	Status string
	Turn   int

	// actual players (not the ones from the spec)
	Players []string
}

func gameStatusFromResponse(id GameID, resp gameStatusResponse) (*GameStatus, error) {
	var visibility string

	sp := GameSpec{
		Description:   resp.Teaser,
		Pace:          resp.Pace,
		AntsPerPlayer: resp.NbAntPerPlayer,
		InitialEnergy: resp.InitialEnergy,
		InitialAcid:   resp.InitialAcid,
	}

	if err := json.Unmarshal(resp.Visibility, &visibility); err == nil {
		sp.Public = visibility == "public"
	} else {
		sp.Public = false

		if err := json.Unmarshal(resp.Visibility, &sp.Players); err != nil {
			return nil, err
		}
	}

	return &GameStatus{
		Game: Game{
			Identifier:   id,
			CreationDate: resp.CreationDate,
			Creator:      resp.Creator,
			Teaser:       resp.Teaser,
			Spec:         &sp,
		},

		Score:   resp.Score,
		Status:  resp.Status.Status,
		Turn:    resp.Turn,
		Players: resp.Players,
	}, nil
}
