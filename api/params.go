package api

// This file describes parameter structs used to send things to the server.
// They're constructed by the client in  `api/client.go` and serialized by the
// low-level HTTP(S) client in `api/io.go`.

// NoParams is used when an API call doesn't need any parameter
type NoParams struct{}

// GameIDParams represent a simple parameters struct for when we need to send a
// game id
type GameIDParams struct {
	ID GameID `url:"id"`
}

// GameSpecParams represents the parameters struct for when we need to send the
// specs for a game (e.g. to create it)
type GameSpecParams struct {
	Users           string `url:"users"`
	Teaser          string `url:"teaser"`
	Pace            int    `url:"pace"`
	NbTurn          int    `url:"nb_turn"`
	NbAntPerPlayer  int    `url:"nb_ant_per_player"`
	NbPlayer        int    `url:"nb_player"`
	MinimalNbPlayer int    `url:"minimal_nb_player"`
	InitialEnergy   int    `url:"initial_energy"`
	InitialAcid     int    `url:"initial_acid"`
}

// UserCredentialsParams represents the parameters struct for the user/password
// params
type UserCredentialsParams struct {
	Login    string `url:"login"`
	Password string `url:"password"`
}

// PlayParams represents the parameters struct needed to play during a turn
type PlayParams struct {
	ID   GameID `url:"id"`
	Cmds string `url:"cmds"`
}

// GenericIDParams is a parameters struct for when we need to send an id
type GenericIDParams struct {
	ID string `url:"id"`
}
