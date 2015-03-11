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

// DumpTo takes a pointer to a struct and dumps the body in it
func (b Body) DumpTo(data interface{}) error {
	if err := b.Error(); err != nil {
		return err
	}

	return json.Unmarshal(*b.Content, data)
}

// JSONString returns a JSON string for this body
func (b Body) JSONString() string {
	return string(*b.Content)
}

// ensureEmptyResponse tries to parse the body as an empty JSON object. It'll
// return an error if .Error() is non-nil
func (b Body) ensureEmptyResponse() (err error) {
	if err = b.Error(); err == nil {
		err = b.DumpTo(&struct{}{})
	}

	return
}

type gameStatusResponse struct {
	Creator        string
	CreationDate   string `json:"creation_date"`
	Teaser         string
	Visibility     string
	NbAntPerPlayer int `json:"nb_ant_per_player"`
	Pace           int
	InitialEnergy  int `json:"initial_energy"`
	InitialAcid    int `json:"initial_acid"`
	Players        []string
	Score          map[string]int
	Status         struct{ Status string }
	Turn           int
}
