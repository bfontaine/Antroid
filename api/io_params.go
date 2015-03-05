package api

// NoParams is used when an API call doesn't need any parameter
type NoParams struct{}

// GameIDParams represent a simple parameters struct for when we need to send a
// game id
type GameIDParams struct {
	Id GameID `url:"id"`
}

// GameSpecParams represents the parameters struct for when we need to send the
// specs for a game (e.g. to create it)
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

// UserCredentialsParams represents the parameters struct for the user/password
// params
type UserCredentialsParams struct {
	Login    string `url:"login"`
	Password string `url:"password"`
}

// PlayParams represents the parameters struct needed to play during a turn
type PlayParams struct {
	Id   GameID `url:"id"`
	Cmds string `url:"cmds"`
}

// GenericIDParams is a parameters struct for when we need to send an id
type GenericIDParams struct {
	Id string `url:"id"`
}
