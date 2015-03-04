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
	Response ApiInfo
}

type gamesResponse struct {
	Status   string
	Response struct {
		// FIXME we don't know the format returned by the API right now
		Games []interface{}
	}
}

type Body struct {
	Content *goreq.Body
	err     error
}

func (b Body) IsEmpty() bool {
	return b.Content == nil
}

func (b Body) Error() (err error) {
	if b.err != nil {
		err = b.err
	}

	if b.IsEmpty() {
		err = ErrEmptyBody
	}

	return
}

func (b Body) FromJsonTo(target interface{}) (err error) {
	err = b.Error()

	if err == nil {
		err = b.Content.FromJsonTo(&target)
	}

	return
}

func (b Body) Close() (err error) {
	if !b.IsEmpty() {
		err = b.Content.Close()
	}

	return
}
