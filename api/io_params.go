package api

// TODO we might be able to remove this package and use annoted version of
// other structs (i.e. add `url:"..."` for each field)

// used when an API call doesn't need any parameter
type NoParams struct{}

// game id params
type GameIdParams struct {
	Id GameId `url: "id"`
}

// game spec params
type GameSpecParams struct {
	Users             string `url:"users"`
	Teaser            string `url:"teaser"`
	Pace              int    `url:"pace"`
	Nb_turn           int    `url:"nb_turn"`
	Nb_ant_per_player int    `url:"nb_ant_per_player"`
	Nb_player         int    `url:"nb_player"`
	Minimal_nb_player int    `url:"minimal_nb_player"`
	Initial_energy    int    `url:"initial_energy"`
	Initial_acid      int    `url:"initial_acid"`
}

// user/password params
type UserCredentialsParams struct {
	Login    string `url:"login"`
	Password string `url:"password"`
}

type PlayParams struct {
	Id   string `url:"id"`
	Cmds string `url:"cmds"`
}

type GenericIdParams struct {
	Id string `url:"id"`
}
