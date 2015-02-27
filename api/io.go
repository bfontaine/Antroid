package api

import (
	"errors"
	"fmt"
	"github.com/franela/goreq"
	"github.com/google/go-querystring/query"
	"strings"
)

/*
Low-level errors that can be returned by the HTTP API.
See the spec: http://yann.regis-gianas.org/antroid/html/api?version=0
*/
var (
	ErrUnknownUser       = errors.New("Unknown user")
	ErrInvalidArgument   = errors.New("Invalid argument")
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

	ErrEmptyBody      = errors.New("Unexpected empty response body")
	ErrNotImplemented = errors.New("Not implemented")

	Err4XX = errors.New("Client error")
	Err5XX = errors.New("Server error")
)

/* The base URL of all API calls */
var BASE_URL = "https://yann.regis-gianas.org/antroid"

/* The API version we support */
var API_VERSION = "0"

/* The User-Agent header we use in all requests */
var USER_AGENT = "Antroid w/ Go, Cailloux&Fontaine&Galichet&Sagot"

// An HTTP client for the API server
type httclient struct {
	UserAgent string

	baseUrl    string
	apiVersion string

	// we use a dead simple implementation here, just storing the cookie as
	// returned by the server and sending it back. We don't check the path nor
	// the protocol nor the expiration date.
	authCookie string
}

// Create a new HTTP client.
func NewHTTClient() httclient {
	return httclient{
		UserAgent:  USER_AGENT,
		baseUrl:    BASE_URL,
		apiVersion: API_VERSION,
		authCookie: "",
	}
}

// Return an absolute URL for a given call.
func (h *httclient) makeApiUrl(call string) string {
	return fmt.Sprintf("%s/%s%s", h.baseUrl, h.apiVersion, string(call))
}

func (h *httclient) setCookieFromHeader(header string) {
	parts := strings.SplitN(header, ";", 2)
	if len(parts) == 0 {
		return
	}

	h.authCookie = parts[0]
}

// Return the appropriate error for a given HTTP code
func (h *httclient) getError(code int) error {
	switch code / 100 {
	case 4:
		return Err4XX
	case 5:
		return Err5XX
	}

	return nil
}

// Make an HTTP call to the remote server and return its response body.
// Don't forget to close it if it's not nil.
func (h *httclient) call(method, call string, data interface{}) (*goreq.Body, error) {

	req := goreq.Request{
		Uri:       h.makeApiUrl(call),
		Method:    string(method),
		Accept:    "application/json",
		UserAgent: h.UserAgent,
		// the server uses a self-signed certificate
		Insecure: true,
		//ShowDebug: true,
	}

	if method == "GET" {
		// goreq will encode everything for us
		req.QueryString = data
	} else {

		// we need to encode our values because the server doesn't accept JSON in
		// requests.
		values, err := query.Values(data)

		if err != nil {
			return nil, err
		}

		queryString := values.Encode()
		req.ContentType = "application/x-www-form-urlencoded"
		req.Body = queryString
	}

	req.AddHeader("Cookie", h.authCookie)

	res, err := req.Do()

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return res.Body, h.getError(res.StatusCode)
	}

	if cookieHeader := res.Header.Get("Set-Cookie"); cookieHeader != "" {
		h.setCookieFromHeader(cookieHeader)
	}

	return res.Body, nil
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
func (h *httclient) CallApi() (*goreq.Body, error) {
	return h.call(get, "/api", NoParams{})
}

// Perform a call to /auth.
func (h *httclient) CallAuth(params UserCredentialsParams) (*goreq.Body, error) {
	return h.call(post, "/auth", params)
}

// Perform a call to /create.
func (h *httclient) CallCreate(params GameSpecParams) (*goreq.Body, error) {
	return h.call(get, "/create", params)
}

// Perform a call to /destroy.
func (h *httclient) CallDestroy(params GameIdParams) (*goreq.Body, error) {
	return h.call(get, "/destroy", params)
}

// Perform a call to /games.
func (h *httclient) CallGames() (*goreq.Body, error) {
	return h.call(get, "/games", NoParams{})
}

// Perform a call to /join.
func (h *httclient) CallJoin(params GameIdParams) (*goreq.Body, error) {
	return h.call(get, "/join", params)
}

// Perform a call to /log.
func (h *httclient) CallLog(params GameIdParams) (*goreq.Body, error) {
	return h.call(get, "/log", params)
}

// Perform a call to /logout.
func (h *httclient) CallLogout() (*goreq.Body, error) {
	return h.call(get, "/logout", NoParams{})
}

// Perform a call to /play.
func (h *httclient) CallPlay(params PlayParams) (*goreq.Body, error) {
	return h.call(get, "/play", params)
}

// Perform a call to /register.
func (h *httclient) CallRegister(params UserCredentialsParams) (*goreq.Body, error) {
	return h.call(post, "/register", params)
}

// Perform a call to /shutdown.
func (h *httclient) CallShutdown(params GenericIdParams) (*goreq.Body, error) {
	return h.call(get, "/shutdown", params)
}

// Perform a call to /status.
func (h *httclient) CallStatus(params GameIdParams) (*goreq.Body, error) {
	return h.call(get, "/status", params)
}

// Perform a call to /whoami.
func (h *httclient) CallWhoAmI() (*goreq.Body, error) {
	return h.call(get, "/whoami", NoParams{})
}
