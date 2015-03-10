package api

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

type GameStatus struct {
	Game

	Score  map[string]int
	Status string
	Turn   int

	// actual players (not the ones from the spec)
	Players []string
}

func gameStatusFromResponse(id GameID, resp gameStatusResponse) (gs *GameStatus) {
	sp := GameSpec{
		Public:        (resp.Visibility == "public"),
		Description:   resp.Teaser,
		Pace:          resp.Pace,
		AntsPerPlayer: resp.NbAntPerPlayer,
		InitialEnergy: resp.InitialEnergy,
		InitialAcid:   resp.InitialAcid,
	}

	gs = &GameStatus{
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
	}

	return
}
