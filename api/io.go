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

var (
	/* The base URL of all API calls */
	BASE_URL = "https://yann.regis-gianas.org/antroid"
	/* The API version we support */
	API_VERSION = "0"
	/* The User-Agent header we use in all requests */
	USER_AGENT = "Antroid w/ Go, Cailloux&Fontaine&Galichet&Sagot"
)

type apiCall string

var (
	CALL_API_INFO = apiCall("/api")
	CALL_AUTH     = apiCall("/auth")
	CALL_CREATE   = apiCall("/create")
	CALL_DESTROY  = apiCall("/destroy")
	CALL_GAMES    = apiCall("/games")
	CALL_JOIN     = apiCall("/join")
	CALL_LOG      = apiCall("/log")
	CALL_LOGOUT   = apiCall("/logout")
	CALL_PLAY     = apiCall("/play")
	CALL_REGISTER = apiCall("/register")
	CALL_SHUTDOWN = apiCall("/shutdown")
	CALL_STATUS   = apiCall("/status")
	CALL_WHOAMI   = apiCall("/whoami")
)

type httpVerb string

var (
	GET  = httpVerb("GET")
	POST = httpVerb("POST")
)

type httclient struct {
	baseUrl    string
	apiVersion string
	userAgent  string
	cookies    cookiejar.Jar
}

func NewHTTClient() httclient {
	return httclient{
		baseUrl:    BASE_URL,
		apiVersion: API_VERSION,
		userAgent:  USER_AGENT,
	}
}

// Return an absolute URL for a given call.
func (h *httclient) MakeApiUrl(call apiCall) string {
	return fmt.Sprintf("%s/%s%s", h.baseUrl, h.apiVersion, string(call))
}

// used when an API call doesn't need any parameter
var NoParams = struct{}{}

// Low-level call
func (h *httclient) Call(method httpVerb, call apiCall, body interface{}) (string, error) {
	res, err := goreq.Request{
		Uri:       h.MakeApiUrl(call),
		Method:    string(method),
		Accept:    "application/json",
		UserAgent: h.userAgent,
		Body:      body,
		// unfortunately the server uses a broken certificate right now
		Insecure: true,
	}.Do()
	// TODO add cookies

	if err != nil {
		return "", err
	}

	// TODO set cookies http://golang.org/pkg/net/http/cookiejar/

	defer res.Body.Close()

	return res.Body.ToString()
}
