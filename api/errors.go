package api

import "errors"

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
	ErrInvalidLogin      = errors.New("Invalid login")

	ErrEmptyBody      = errors.New("Unexpected empty response body")
	ErrUnknown        = errors.New("Unknown error")
	ErrUnknownCode    = errors.New("Unknown error code")
	ErrNotImplemented = errors.New("Not implemented")

	Err4XX = errors.New("Client error")
	Err5XX = errors.New("Server error")
)

// See the API spec
var errorCodes = map[int]error{
	202165063:  ErrUnknownUser,
	285625267:  ErrWrongCmd,
	306276868:  ErrAlreadyJoined,
	318351321:  ErrGameNotPlaying,
	332299703:  ErrUserAlreadyExists,
	415302510:  ErrNotLogged,
	591857505:  ErrWrongAnt,
	598240942:  ErrGameNotOver,
	621433138:  ErrInvalidLogin,
	683983482:  ErrMustJoin,
	761507830:  ErrNoMoreSlot,
	796193025:  ErrWrongGame,
	995492770:  ErrInvalidArgument,
	1032614003: ErrNoPerm,
}

func errorForCode(code int) (err error) {
	err, ok := errorCodes[code]
	if !ok {
		err = ErrUnknownCode
	}

	return
}
