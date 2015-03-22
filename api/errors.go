package api

// This files describes all errors that can be returned by functions and
// methods under `api/*`. Most of them are just transcriptions of errors from
// the remote API server
//
// See also http://yann.regis-gianas.org/antroid/html/api?version=0

import "errors"

var (
	// Remote API errors

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
	ErrInvalidLogin      = errors.New("Invalid login")

	// Our errors

	ErrEmptyBody      = errors.New("Unexpected empty response body")
	ErrNotImplemented = errors.New("Not implemented")

	// This error is returned when the remote server returns an error with a
	// status we don't understand (nor "completed" nor "errors"). This
	// shouldn't happen in practice.
	ErrUnknown = errors.New("Unknown error")

	// This error is returned when the remote server gives us an unknown error
	// code.
	ErrUnknownCode = errors.New("Unknown error code")

	// HTTP errors

	Err4XX = errors.New("Client error")
	Err5XX = errors.New("Server error")
)

// See the API spec:
//   http://yann.regis-gianas.org/antroid/html/api?version=0
var errorCodes = map[int]error{
	502441794:  ErrUnknownUser,       // UNKNOWN_USER
	4313039:    ErrWrongCmd,          // INVALID_COMMAND
	443034632:  ErrAlreadyJoined,     // ALREADY_JOINED
	357629463:  ErrGameNotPlaying,    // GAME_IS_NOT_PLAYING
	334269347:  ErrUserAlreadyExists, // USER_ALREADY_EXISTS
	32403037:   ErrNotLogged,         // MUST_BE_LOGGED
	523177601:  ErrMustJoin,          // MUST_JOIN_FIRST
	973268714:  ErrWrongAnt,          // INVALID_ANT_IDENTIFIER
	808783211:  ErrGameNotOver,       // GAME_IS_NOT_OVER
	114981602:  ErrInvalidLogin,      // INVALID_LOGIN
	1001223883: ErrNoMoreSlot,        // NO_MORE_SLOT
	513589683:  ErrWrongGame,         // INVALID_GAME_IDENTIFIER
	677388348:  ErrInvalidArgument,   // INVALID_ARGUMENT
	842913797:  ErrNoPerm,            // NO_PERMISSION
}

// errorForCode returns a Go error from an error code. If we don't know the
// code, it'll return ErrUnknownCode.
func errorForCode(code int) (err error) {
	err, ok := errorCodes[code]
	if !ok {
		err = ErrUnknownCode
	}

	return
}
