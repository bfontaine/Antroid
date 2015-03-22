package api

// This file describes structs used to represent a game and its current status

import "encoding/json"

// GameID is a game id
type GameID string

// Game represents a game
type Game struct {
	Identifier   GameID
	CreationDate string `json:"creation_date"`
	Creator      string
	Teaser       string

	// The spec contains all the parameters such as the pace, number of ants
	// per player, initial energy level, etc.
	Spec *GameSpec
}

// A GameStatus represents the current status of a game
type GameStatus struct {
	// It inherits all fields of a game
	Game

	// The current scoreboard (map username => score)
	Score map[string]int
	// The game status, e.g.: "over", "playing". We don't know if there are
	// more of them, which is why we use a string here instead of an int
	// constant.
	Status string
	// The current turn number
	Turn int

	// actual players (not the ones from the spec)
	Players []string
}

// A GameLog is a log from a game. We haven't implemented it yet.
type GameLog struct{}

// This is an internal helpers which constructs a GameStatus struct from a
// response from the remote server
func gameStatusFromResponse(id GameID, resp gameStatusResponse) (*GameStatus, error) {
	var visibility string

	// first, fill the spec
	sp := GameSpec{
		Description:   resp.Teaser,
		Pace:          resp.Pace,
		AntsPerPlayer: resp.NbAntPerPlayer,
		InitialEnergy: resp.InitialEnergy,
		InitialAcid:   resp.InitialAcid,
	}

	// The API returns either "public" or an array of players. We first try to
	// parse the visibility as a string, ...
	if err := json.Unmarshal(resp.Visibility, &visibility); err == nil {
		sp.Public = visibility == "public"
	} else {
		sp.Public = false

		// ...and fallback on an array.
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
