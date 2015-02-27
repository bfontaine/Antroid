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

func (cl *Client) getUserCredentialsParams() UserCredentialsParams {
	return UserCredentialsParams{
		Login:    cl.username,
		Password: cl.password,
	}
}

// Test if the client is authenticated.
func (cl *Client) Authenticated() bool {
	return cl.authenticated
}

// Get some info about the API.
func (cl *Client) ApiInfo() (ApiInfo, error) {
	body, err := cl.http.CallApi()

	if err != nil {
		return ApiInfo{}, err
	}

	if body == nil {
		return ApiInfo{}, ErrEmptyBody
	}

	defer body.Close()

	var resp apiInfoResponse

	body.FromJsonTo(&resp)

	if resp.Status != "completed" {
		return ApiInfo{}, ErrUnknown
	}

	return resp.Response, nil
}

// Register some credentials for this client
func (cl *Client) RegisterWithCredentials(username, password string) error {
	cl.username = username
	cl.password = password

	body, err := cl.http.CallRegister(cl.getUserCredentialsParams())

	if err != nil {
		return nil
	}

	var resp registerResponse
	body.FromJsonTo(&resp)

	body.Close()

	if resp.Status == "error" {
		return errorForCode(resp.Response.Error_code)
	}

	return nil
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
	body, err := cl.http.CallAuth(cl.getUserCredentialsParams())

	cl.authenticated = (err == nil)

	// TODO check the body
	if body != nil {
		body.Close()
	}

	return err
}

// Logout the client.
// If the client wasn't already authenticated the method returns without
// failing.
func (cl *Client) Logout() error {
	if !cl.authenticated {
		return nil
	}

	// not tested
	_, err := cl.http.CallLogout()

	return err
}

// Create a new game.
func (cl *Client) CreateGame(gs *GameSpec) (Game, error) {
	_, err := cl.http.CallCreate(gs.toParams())

	if err != nil {
		return Game{}, err
	}

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
	return nil
}

// Destroy a game given its identifier.
func (cl *Client) DestroyGameIndentifier(id GameId) error {
	// TODO
	return ErrNotImplemented
}

// List all visible games.
func (cl *Client) ListGames() ([]Game, error) {
	_, err := cl.http.CallGames()
	if err != nil {
		return nil, err
	}

	// TODO

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

// Shutdown a server (need to be root). We don't know what's this id for
func (cl *Client) ShutdownIdentifier(id string) error {
	if !cl.authenticated {
		return ErrNoPerm
	}

	_, err := cl.http.CallShutdown(GenericIdParams{Id: id})

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

// number of characters we need to skip in /whoami's response to get our
// username.
var whoAmILoginSlice = len("logged as ")

// Check the client's status on the server-side and return it.
func (cl *Client) WhoAmI() (string, error) {
	body, err := cl.http.CallWhoAmI()

	if err != nil {
		return "", err
	}

	if body == nil {
		return "", ErrEmptyBody
	}

	defer body.Close()

	var resp whoAmIResponse

	body.FromJsonTo(&resp)

	cl.authenticated = (resp.Response.Status != "not_logged")

	if !cl.authenticated {
		return "", ErrNotLogged
	}

	return resp.Response.Status[whoAmILoginSlice:], nil
}
