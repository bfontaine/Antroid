package api

// This file describes an API client, which is a high-level object used to
// access the remote API. See `api/io.go` for the low-level I/O stuff.
//
// See `api/errors.go` for an overview of all errors that can be returned by
// these methods.
//
// All methods do roughly the same thing:
// 1. build the API parameters from its own parameters
// 2. call the low-level HTTP(S) client
// 3. check if it returned an error, and if so return it
// 4. unparse its output in a structure
// 5. return the structure
//
// Most of the code here is just a high-level wrapper around the low-level
// HTTP(S) client with checks for errors everywhere.

// Client is an API client
type Client struct {
	// username and password used for authentication
	username, password string
	// we remember if we're logged on the server or not
	authenticated bool
	// the low-level HTTP(S) client
	http *Httclient
	// if `debug` is true, we'll be more verbose
	debug bool
}

// NewClient creates and returns a new API client.
func NewClient() *Client {
	return &Client{
		http: NewHTTClient(),
	}
}

// SetDebug sets the debug flag
func (cl *Client) SetDebug(debug bool) {
	cl.debug = debug
	cl.http.debug = debug
}

// getUserCredentialsParams returns the client's credentials (username and
// password) as an UserCredentialsParams struct, which can then be passed to
// the low-level HTTP(S) client.
func (cl *Client) getUserCredentialsParams() UserCredentialsParams {
	return UserCredentialsParams{
		Login:    cl.username,
		Password: cl.password,
	}
}

// Authenticated tests if the client is authenticated. This is an exported
// getter for our private `authenticated` field.
func (cl *Client) Authenticated() bool {
	return cl.authenticated
}

// APIInfo returns some info about the API. See `api/api_info.go` for the
// returned struct.
func (cl *Client) APIInfo() (info APIInfo, err error) {
	body := cl.http.CallAPI()

	if err = body.Error(); err == nil {
		// if we didn't get an error on parsing and the server didn't send an
		// error code, extract the JSON in an APIInfo struct.
		err = body.DumpTo(&info)
	}

	return
}

// RegisterWithCredentials registers some credentials for this client, creating
// an account on the remote server.
func (cl *Client) RegisterWithCredentials(username, password string) error {
	cl.username = username
	cl.password = password

	body := cl.http.CallRegister(cl.getUserCredentialsParams())

	// the server shouldn't send back anything except a success status
	return body.ensureEmptyResponse()
}

// LoginWithCredentials authenticates the client with the given credentials.
// If the client was already authenticated and the new username/password are
// the same it's not re-authenticated and the method returns without failing.
// In any other case an API call is made. The returned error can be either nil
// (success) or ErrUnknownUser.
func (cl *Client) LoginWithCredentials(username, password string) error {
	if cl.authenticated {
		if cl.username == username && cl.password == password {
			return nil
		}

		// if we're already authenticated but with different credentials,
		// logout us before re-logging-in with the new credentials. This is not
		// useful on our side but it's to be polite with the server.
		cl.Logout()
	}

	// update our credentials
	cl.username = username
	cl.password = password

	// and login with it
	return cl.Login()
}

// Login anthenticates the client with its own credentials.
func (cl *Client) Login() (err error) {
	body := cl.http.CallAuth(cl.getUserCredentialsParams())

	err = body.Error()

	// we're authenticated if the server didn't respond with an error
	cl.authenticated = (err == nil)

	if err == nil {
		err = body.ensureEmptyResponse()
	}

	return
}

// Logout the client.
// If the client wasn't already authenticated the method returns without
// failing.
func (cl *Client) Logout() (err error) {
	if !cl.authenticated {
		return
	}

	b := cl.http.CallLogout()

	if err = b.Error(); err != nil {
		return
	}

	if err = b.ensureEmptyResponse(); err == nil {
		// we're deconnected only if the server returned a non-error status
		cl.authenticated = false
	}

	return
}

// CreateGame creates a new game and returns it. See `api/game_spec.go` for
// what is a GameSpec struct (spoiler: it describes a game's params).
func (cl *Client) CreateGame(gs *GameSpec) (g *Game, err error) {
	body := cl.http.CallCreate(gs.toParams())

	if err = body.Error(); err != nil {
		return
	}

	// The API returns a game identifier
	var resp struct{ Identifier string }

	if err = body.DumpTo(&resp); err != nil {
		return
	}

	// create a new game from the spec and the identifier
	return &Game{Spec: gs, Identifier: GameID(resp.Identifier)}, nil
}

// DestroyGame destroys a game.
// If the method is successful it'll modify the game in-place and reset its
// identifier.
func (cl *Client) DestroyGame(g *Game) (err error) {
	// this is just a proxy to DestroyGameIdentifier
	if err = cl.DestroyGameIdentifier(g.Identifier); err == nil {
		g.Identifier = ""
	}
	return
}

// DestroyGameIdentifier destroys a game given its identifier.
func (cl *Client) DestroyGameIdentifier(id GameID) (err error) {
	body := cl.http.CallDestroy(GameIDParams{ID: id})

	return body.ensureEmptyResponse()
}

// ListGames lists all visible games.
func (cl *Client) ListGames() (games []Game, err error) {
	body := cl.http.CallGames()

	if err = body.Error(); err != nil {
		return
	}

	// the response contains a list of JSON objects with a `game_description`
	// key that describes a game.
	var resp struct {
		Games []struct {
			GameDescription Game `json:"game_description"`
		}
	}

	if err = body.DumpTo(&resp); err != nil {
		return
	}

	for _, g := range resp.Games {
		games = append(games, g.GameDescription)
	}

	return
}

// JoinGame makes the client join a game
func (cl *Client) JoinGame(g *Game) error {
	// just a proxy to JoinGameIdentifier
	return cl.JoinGameIdentifier(g.Identifier)
}

// JoinGameIdentifier makes the client join a game given its identifier.
func (cl *Client) JoinGameIdentifier(id GameID) error {
	body := cl.http.CallJoin(GameIDParams{ID: id})

	return body.ensureEmptyResponse()
}

// GetGameLog returns a game's log
func (cl *Client) GetGameLog(g *Game) (GameLog, error) {
	// just a proxy to GetGameIdentifierLog
	return cl.GetGameIdentifierLog(g.Identifier)
}

// GetGameIdentifierLog returns a game's log given its identifier
func (cl *Client) GetGameIdentifierLog(id GameID) (gl GameLog, err error) {
	body := cl.http.CallLog(GameIDParams{ID: id})

	if err = body.Error(); err != nil {
		return
	}

	// TODO
	// the log was buggy when we wrote this part so this is still not
	// implemented

	return GameLog{}, ErrNotImplemented
}

// Play plays a game with a list of commands
func (cl *Client) Play(g *Game, cmds Commands) (*Turn, error) {
	// just a proxy to PlayIdentifier
	return cl.PlayIdentifier(g.Identifier, cmds)
}

// PlayIdentifier plays a game with a list of commands, given its identifier
func (cl *Client) PlayIdentifier(id GameID, cmds Commands) (t *Turn, err error) {
	body := cl.http.CallPlay(PlayParams{ID: id, Cmds: cmds.String()})

	if err = body.Error(); err != nil {
		return
	}

	var resp playResponse

	if err = body.DumpTo(&resp); err == nil {
		t, err = resp.getTurn()
	}

	return
}

// ShutdownIdentifier shutdowns a server (need to be root). We don't know
// what's this id for.
func (cl *Client) ShutdownIdentifier(id string) error {
	// it's useless to implement that since we don't have root credentials.
	return ErrNotImplemented
}

// GetGameStatus returns a game's status
func (cl *Client) GetGameStatus(g *Game) (*GameStatus, error) {
	// just a proxy to GetGameIdentifierStatus
	return cl.GetGameIdentifierStatus(g.Identifier)
}

// GetGameIdentifierStatus gets a game's status, given its identifier
func (cl *Client) GetGameIdentifierStatus(id GameID) (gs *GameStatus, err error) {
	body := cl.http.CallStatus(GameIDParams{ID: id})

	if err = body.Error(); err != nil {
		return
	}

	var resp struct{ Status gameStatusResponse }

	if err = body.DumpTo(&resp); err != nil {
		return
	}

	gs, err = gameStatusFromResponse(id, resp.Status)

	return
}

// number of characters we need to skip in /whoami's response to get our
// username.
var whoAmILoginSlice = len("logged as ")

// WhoAmI checks the client's status on the server-side and return it.
func (cl *Client) WhoAmI() (s string, err error) {
	body := cl.http.CallWhoAmI()

	if err = body.Error(); err != nil {
		return
	}

	var resp struct{ Status string }

	if err = body.DumpTo(&resp); err != nil {
		return
	}

	st := resp.Status

	// we're authenticated if the server didn't told use we aren't
	cl.authenticated = (st != "" && st != "not_logged")

	if !cl.authenticated {
		err = ErrNotLogged
		return
	}

	// extract our username by skipping "logged as " at the beginning of the
	// string (TODO we could use TrimPrefix here)
	return st[whoAmILoginSlice:], nil
}
