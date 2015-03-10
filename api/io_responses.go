package api

import (
	"encoding/json"
)

/*
   This file defines structures used to unmarshal JSON responses from the API.
*/

type baseResponse struct {
	Status   string
	Response json.RawMessage
}

type errorResponse struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_msg"`
}

func (e errorResponse) Error() error {
	return errorForCode(e.Code)
}

// Body is a body response from the API
type Body struct {
	Content    *json.RawMessage
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

func (b Body) DumpTo(data interface{}) error {
	if err := b.Error(); err != nil {
		return err
	}
	if b.IsEmpty() {
		return ErrEmptyBody
	}

	return json.Unmarshal(*b.Content, data)
}

func (b Body) JSONString() string {
	return string(*b.Content)
}

func (b Body) ensureEmptyResponse() error {
	return b.DumpTo(&struct{}{})
}
