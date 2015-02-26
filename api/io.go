package api

import (
	"errors"
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
