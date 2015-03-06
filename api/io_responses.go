package api

import (
	"github.com/franela/goreq"
)

/*
   This file defines structures used to unmarshal JSON responses from the API.
*/

type simpleResponse struct {
	Status   string
	Response struct {
		Error_code int
		Error_msg  string

		// used for game creation
		Identifier string

		// used for whoami calls
		Status string
	}
}

func (r simpleResponse) IsError() bool {
	return r.Status == "error"
}

func (r simpleResponse) Error() error {
	if !r.IsError() {
		return nil
	}
	return errorForCode(r.Response.Error_code)
}

type apiInfoResponse struct {
	Status   string
	Response APIInfo
}

type gamesResponse struct {
	Status   string
	Response struct {
		// FIXME we don't know the format returned by the API right now
		Games []interface{}
	}
}

// Body is a body response from the API
type Body struct {
	Content    *goreq.Body
	StatusCode int
	err        error
}

// IsEmpty returns true if the body response is empty
func (b Body) IsEmpty() bool {
	return b.Content == nil
}

// Error returns any error with this body
func (b Body) Error() (err error) {
	if b.err != nil {
		return b.err
	}

	if b.IsEmpty() {
		err = ErrEmptyBody
	}

	return
}

// FromJSONTo assumes the body contains JSON and dumps it in the given struct
func (b Body) FromJSONTo(target interface{}) (err error) {
	err = b.Error()

	if err == nil {
		err = b.Content.FromJsonTo(&target)
	}

	return
}

// Close closes the body if it's not empty
func (b Body) Close() (err error) {
	if !b.IsEmpty() {
		err = b.Content.Close()
	}

	return
}
