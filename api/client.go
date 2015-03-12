package api

// Client is an API client
type Client struct {
	username      string
	password      string
	authenticated bool
	http          *Httclient

	debug bool
}

// NewClient creates and returns a new API client.
func NewClient() (*Client, error) {
	return &Client{
		http: NewHTTClient(),
	}, nil
}

// SetDebug sets the debug flag
func (cl *Client) SetDebug(debug bool) {
	cl.debug = debug
	cl.http.debug = debug
}

func (cl *Client) getUserCredentialsParams() UserCredentialsParams {
	return UserCredentialsParams{
		Login:    cl.username,
		Password: cl.password,
	}
}

// Authenticated tests if the client is authenticated.
func (cl *Client) Authenticated() bool {
	return cl.authenticated
}

// APIInfo returns some info about the API.
func (cl *Client) APIInfo() (info APIInfo, err error) {
	body := cl.http.CallAPI()

	if err = body.Error(); err == nil {
		err = body.DumpTo(&info)
	}

	return
}

// RegisterWithCredentials registers some credentials for this client
func (cl *Client) RegisterWithCredentials(username, password string) error {
	cl.username = username
	cl.password = password

	body := cl.http.CallRegister(cl.getUserCredentialsParams())

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

		cl.Logout()
	}

	cl.username = username
	cl.password = password

	return cl.Login()
}

// Login anthenticates the client with its own credentials.
func (cl *Client) Login() (err error) {
	body := cl.http.CallAuth(cl.getUserCredentialsParams())

	err = body.Error()

	cl.authenticated = (err == nil)

	if err != nil {
		return
	}

	return body.ensureEmptyResponse()
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
		cl.authenticated = false
	}

	return
}

// CreateGame creates a new game and returns it.
func (cl *Client) CreateGame(gs *GameSpec) (g Game, err error) {
	body := cl.http.CallCreate(gs.toParams())

	if err = body.Error(); err != nil {
		return
	}

	var resp struct{ Identifier string }

	if err = body.DumpTo(&resp); err != nil {
		return
	}

	g.Spec = gs
	g.Identifier = GameID(resp.Identifier)

	return
}

// DestroyGame destroys a game.
// If the method is successful it'll modify the game in-place and reset its
// identifier.
func (cl *Client) DestroyGame(g *Game) error {
	if err := cl.DestroyGameIdentifier(g.Identifier); err != nil {
		return err
	}
	g.Identifier = ""
	return nil
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

	var resp struct {
		Games []struct{ Game_description Game }
	}

	if err = body.DumpTo(&resp); err != nil {
		return
	}

	for _, g := range resp.Games {
		games = append(games, g.Game_description)
	}

	return
}

// JoinGame makes the client join a game
func (cl *Client) JoinGame(g *Game) error {
	return cl.JoinGameIdentifier(g.Identifier)
}

// JoinGameIdentifier makes the client join a game given its identifier.
func (cl *Client) JoinGameIdentifier(id GameID) error {
	body := cl.http.CallJoin(GameIDParams{ID: id})

	return body.ensureEmptyResponse()
}

// GetGameLog returns a game's log
func (cl *Client) GetGameLog(g *Game) (GameLog, error) {
	return cl.GetGameIdentifierLog(g.Identifier)
}

// GetGameIdentifierLog returns a game's log given its identifier
func (cl *Client) GetGameIdentifierLog(id GameID) (gl GameLog, err error) {
	body := cl.http.CallLog(GameIDParams{ID: id})

	if err = body.Error(); err != nil {
		return
	}

	// TODO

	return GameLog{}, ErrNotImplemented
}

// Play a game with a list of commands
func (cl *Client) Play(g *Game, cmds Commands) (*Turn, error) {
	return cl.PlayIdentifier(g.Identifier, cmds)
}

// Play a game with a list of commands, given its identifier
func (cl *Client) PlayIdentifier(id GameID, cmds Commands) (t *Turn, err error) {
	body := cl.http.CallPlay(PlayParams{ID: id, Cmds: cmds.String()})

	if err = body.Error(); err != nil {
		return
	}

	var resp playResponse

	if err = body.DumpTo(&resp); err != nil {
		return
	}

	return resp.getTurn()
}

// ShutdownIdentifier shutdowns a server (need to be root). We don't know
// what's this id for.
func (cl *Client) ShutdownIdentifier(id string) error {
	return ErrNotImplemented
}

// GetGameStatus returns a game's status
// Note: the spec is unclear on the returned JSON so we can't set a return type
// now.
func (cl *Client) GetGameStatus(g *Game) (*GameStatus, error) {
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

	gs = gameStatusFromResponse(id, resp.Status)

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

	cl.authenticated = (st != "" && st != "not_logged")

	if !cl.authenticated {
		err = ErrNotLogged
		return
	}

	s = st[whoAmILoginSlice:]
	return
}
