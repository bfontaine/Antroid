package api

type Client struct {
	username      string
	password      string
	authenticated bool
	http          httclient
}

// Create a new API client.
func NewClient() (Client, error) {
	return Client{
		http: NewHTTClient(),
	}, nil
}

// Test if the client is authenticated.
func (cl *Client) Authenticated() bool {
	return cl.authenticated
}

// Get some info about the API.
func (cl *Client) ApiInfo() (ApiInfo, error) {
	// TODO
	return ApiInfo{}, ErrNotImplemented
}

// Register some credentials for this client
func (cl *Client) RegisterWithCredentials(username, password string) error {
	cl.username = username
	cl.password = password

	return ErrNotImplemented
}

// Authenticate the client with the given credentials.
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

// Anthenticate the client with its own credentials.
func (cl *Client) Login() error {
	return ErrNotImplemented
}

// Logout the client.
// If the client wasn't already authenticated the method returns without
// failing.
func (cl *Client) Logout() error {
	if !cl.authenticated {
		return nil
	}

	// not tested
	_, err := cl.http.Call(POST, CALL_LOGOUT, NoParams)

	return err
}

// Create a new game.
func (cl *Client) CreateGame(g *GameSpec) (Game, error) {
	// TODO
	return Game{}, ErrNotImplemented
}

// Destroy a game.
// If the method is successful it'll modify the game in-place and reset its
// identifier.
func (cl *Client) DestroyGame(g *Game) error {
	if err := cl.DestroyGameIndentifier(g.Identifier); err != nil {
		return err
	}
	g.Identifier = ""
	return ErrNotImplemented
}

// Destroy a game given its identifier.
func (cl *Client) DestroyGameIndentifier(id GameId) error {
	// TODO
	return ErrNotImplemented
}

// List all visible games.
func (cl *Client) ListGames() ([]Game, error) {
	res, err := cl.http.Call(GET, CALL_GAMES, NoParams)
	if err != nil {
		return nil, err
	}

	// TODO

	// trick go to avoid "'res' variable not used" error
	if res == "" {
	}

	return nil, ErrNotImplemented
}

// Join a game
func (cl *Client) JoinGame(g *Game) error {
	return cl.JoinGameIdentifier(g.Identifier)
}

// Join a game given its identifier.
func (cl *Client) JoinGameIdentifier(id GameId) error {
	return ErrNotImplemented
}

// Get a game's log
func (cl *Client) GetGameLog(g *Game) (GameLog, error) {
	return cl.GetGameIdentifierLog(g.Identifier)
}

// Get a game's log given its identifier
func (cl *Client) GetGameIdentifierLog(id GameId) (GameLog, error) {
	return GameLog{}, ErrNotImplemented
}

// Play a game with a list of commands
func (cl *Client) Play(g *Game, cmds []*Command) error {
	return ErrNotImplemented
}

// Shutdown a game (need to be root)
func (cl *Client) ShutdownGame(g *Game) error {
	return cl.ShutdownGameIdentifier(g.Identifier)
}

// Shutdown a game given its identifier (need to be root)
func (cl *Client) ShutdownGameIdentifier(id GameId) error {
	if !cl.authenticated {
		return ErrNoPerm
	}

	_, err := cl.http.Call(GET, CALL_SHUTDOWN, GameIdParams{id: id})

	return err
}

// Get a game's status
// Note: the spec is unclear on the returned JSON so we can't set a return type
// now.
func (cl *Client) GetGameStatus(g *Game) error {
	return cl.GetGameIdentifierStatus(g.Identifier)
}

// Get a game's status, given its identifier
func (cl *Client) GetGameIdentifierStatus(id GameId) error {
	return ErrNotImplemented
}

// Check the client's status on the server-side
func (cl *Client) WhoAmI() (string, error) {
	return "", ErrNotImplemented
}
