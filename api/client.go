package api

// Client is an API client
type Client struct {
	username      string
	password      string
	authenticated bool
	http          *Httclient
}

// NewClient creates and returns a new API client.
func NewClient() (*Client, error) {
	return &Client{
		http: NewHTTClient(),
	}, nil
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

	if err := body.Error(); err != nil {
		return err
	}

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

		if err := cl.Logout(); err != nil {
			return err
		}
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

	if err = b.ensureEmptyResponse(); err != nil {
		return
	}

	cl.authenticated = false
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
	if err := cl.DestroyGameIndentifier(g.Identifier); err != nil {
		return err
	}
	g.Identifier = ""
	return nil
}

// DestroyGameIndentifier destroys a game given its identifier.
func (cl *Client) DestroyGameIndentifier(id GameID) error {
	// TODO
	return ErrNotImplemented
}

// ListGames lists all visible games.
func (cl *Client) ListGames() (games []Game, err error) {
	body := cl.http.CallGames()

	if err = body.Error(); err != nil {
		return
	}

	var resp struct{ Games []Game }

	if err = body.DumpTo(&resp); err != nil {
		return
	}

	//return resp.Response.Games, nil
	err = ErrNotImplemented
	return
}

// JoinGame makes the client join a game
func (cl *Client) JoinGame(g *Game) error {
	return cl.JoinGameIdentifier(g.Identifier)
}

// JoinGameIdentifier makes the client join a game given its identifier.
func (cl *Client) JoinGameIdentifier(id GameID) error {
	return ErrNotImplemented
}

// GetGameLog returns a game's log
func (cl *Client) GetGameLog(g *Game) (GameLog, error) {
	return cl.GetGameIdentifierLog(g.Identifier)
}

// GetGameIdentifierLog returns a game's log given its identifier
func (cl *Client) GetGameIdentifierLog(id GameID) (GameLog, error) {
	return GameLog{}, ErrNotImplemented
}

// Play a game with a list of commands
func (cl *Client) Play(g *Game, cmds []*Command) error {
	return ErrNotImplemented
}

// ShutdownIdentifier shutdowns a server (need to be root). We don't know
// what's this id for.
func (cl *Client) ShutdownIdentifier(id string) error {
	if !cl.authenticated {
		return ErrNoPerm
	}

	body := cl.http.CallShutdown(GenericIDParams{ID: id})

	return body.Error()
}

// GetGameStatus returns a game's status
// Note: the spec is unclear on the returned JSON so we can't set a return type
// now.
func (cl *Client) GetGameStatus(g *Game) error {
	return cl.GetGameIdentifierStatus(g.Identifier)
}

// GetGameIdentifierStatus gets a game's status, given its identifier
func (cl *Client) GetGameIdentifierStatus(id GameID) error {
	return ErrNotImplemented
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
