package api

import (
	"errors"
	"fmt"
	"github.com/franela/goreq"
	"net/http/cookiejar"
)

/*
Low-level errors that can be returned by the HTTP API.
See the spec: http://yann.regis-gianas.org/antroid/html/api?version=0
*/
var (
	ErrUnknownUser       = errors.New("Unknown user")
	ErrGameAlreadyExists = errors.New("Game identifier already exists")
	ErrUserAlreadyExists = errors.New("User already exists")
	ErrWrongGame         = errors.New("Invalid game identifier")
	ErrNoPerm            = errors.New("No permission")
	ErrNoMoreSlot        = errors.New("No more slot")
	ErrAlreadyJoined     = errors.New("Already joined")
	ErrGameNotOver       = errors.New("The game is not over")
	ErrWrongAnt          = errors.New("Invalid ant identifier")
	ErrMustJoin          = errors.New("Must join first")
	ErrNotLogged         = errors.New("Must be logged")
	ErrWrongCmd          = errors.New("Invalid command")
	ErrGameNotPlaying    = errors.New("The game is not playing")
	ErrNotImplemented    = errors.New("Not implemented")
)

/* The base URL of all API calls */
var BASE_URL = "https://yann.regis-gianas.org/antroid"

/* The API version we support */
var API_VERSION = "0"

/* The User-Agent header we use in all requests */
var USER_AGENT = "Antroid w/ Go, Cailloux&Fontaine&Galichet&Sagot"

// An HTTP client with a cookiejar
type httclient struct {
	UserAgent string

	baseUrl    string
	apiVersion string
	cookies    cookiejar.Jar
}

// Create a new HTTP client.
func NewHTTClient() httclient {
	// TODO setup cookiejar

	return httclient{
		baseUrl:    BASE_URL,
		apiVersion: API_VERSION,
		UserAgent:  USER_AGENT,
	}
}

// Return an absolute URL for a given call.
func (h *httclient) makeApiUrl(call string) string {
	return fmt.Sprintf("%s/%s%s", h.baseUrl, h.apiVersion, string(call))
}

// Make an HTTP call to the remote server and return its response.
func (h *httclient) call(method, call string, data interface{}) (string, error) {
	req := goreq.Request{
		Uri:       h.makeApiUrl(call),
		Method:    string(method),
		Accept:    "application/json",
		UserAgent: h.UserAgent,
		// the server uses a self-signed certificate
		Insecure: true,
	}

	if method == "GET" {
		req.QueryString = data
	} else {
		req.Body = data
	}

	// TODO add cookies

	res, err := req.Do()

	if err != nil {
		return "", err
	}

	// TODO set cookies http://golang.org/pkg/net/http/cookiejar/

	defer res.Body.Close()

	// TODO extract the error code if there's one
	return res.Body.ToString()
}

var (
	get  = "GET"
	post = "POST"
)

/*
   Each method below perform a call to one endpoint. We expose them instead of
   the generic .call method to be able to type-check the parameters of each
   call.
*/

// Perform a call to /api.
func (h *httclient) CallApi() (string, error) {
	return h.call(get, "/api", NoParams{})
}

// Perform a call to /auth.
func (h *httclient) CallAuth(params UserCredentialsParams) (string, error) {
	return h.call(post, "/auth", params)
}

// Perform a call to /create.
func (h *httclient) CallCreate(params GameSpecParams) (string, error) {
	return h.call(get, "/create", params)
}

// Perform a call to /destroy.
func (h *httclient) CallDestroy(params GameIdParams) (string, error) {
	return h.call(get, "/destroy", params)
}

// Perform a call to /games.
func (h *httclient) CallGames() (string, error) {
	return h.call(get, "/games", NoParams{})
}

// Perform a call to /join.
func (h *httclient) CallJoin(params GameIdParams) (string, error) {
	return h.call(get, "/join", params)
}

// Perform a call to /log.
func (h *httclient) CallLog(params GameIdParams) (string, error) {
	return h.call(get, "/log", params)
}

// Perform a call to /logout.
func (h *httclient) CallLogout() (string, error) {
	return h.call(get, "/logout", NoParams{})
}

// Perform a call to /play.
func (h *httclient) CallPlay(params PlayParams) (string, error) {
	return h.call(get, "/play", params)
}

// Perform a call to /register.
func (h *httclient) CallRegister(params UserCredentialsParams) (string, error) {
	return h.call(post, "/register", params)
}

// Perform a call to /shutdown.
func (h *httclient) CallShutdown(params GenericIdParams) (string, error) {
	return h.call(get, "/shutdown", params)
}

// Perform a call to /status.
func (h *httclient) CallStatus(params GameIdParams) (string, error) {
	return h.call(get, "/status", params)
}

// Perform a call to /whoami.
func (h *httclient) CallWhoAmI() (string, error) {
	return h.call(get, "/whoami", NoParams{})
}
