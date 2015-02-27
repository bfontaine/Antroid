package api

// used when an API call doesn't need any parameter
type NoParams struct{}

// game id params
type GameIdParams struct {
	id GameId
}

// game spec params
type GameSpecParams struct {
	users             string
	teaser            string
	pace              int
	nb_turn           int
	nb_ant_per_player int
	nb_player         int
	minimal_nb_player int
	initial_energy    int
	initial_acid      int
}

// user/password params
type UserCredentialsParams struct {
	user     string
	password string
}

type PlayParams struct {
	id   string
	cmds string
}

type GenericIdParams struct {
	id string
}
