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

	ErrEmptyBody      = errors.New("Unexpected empty response body")
	ErrUnknown        = errors.New("Unknown error")
	ErrUnknownCode    = errors.New("Unknown error code")
	ErrNotImplemented = errors.New("Not implemented")

	Err4XX = errors.New("Client error")
	Err5XX = errors.New("Server error")
)

// See the API spec
var errorCodes = map[int]error{
	61760457:   ErrWrongCmd,
	91411898:   ErrNoMoreSlot,
	351345662:  ErrNoPerm,
	427619750:  ErrWrongGame,
	448649162:  ErrGameNotPlaying,
	565287715:  ErrAlreadyJoined,
	591603053:  ErrInvalidArgument,
	693680202:  ErrUnknownUser,
	813909381:  ErrMustJoin,
	873213279:  ErrUserAlreadyExists,
	942350302:  ErrWrongAnt,
	965395831:  ErrGameNotOver,
	1058501022: ErrNotLogged,
}

func errorForCode(code int) (err error) {
	err, ok := errorCodes[code]
	if !ok {
		err = ErrUnknownCode
	}

	return
}
