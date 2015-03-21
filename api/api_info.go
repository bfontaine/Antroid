package api

// This file describes the struct APIInfo which holds info about the remote
// API. You shouldn't construct these structs yourself, they are returned by
// the API Client (see `api/client.go`).

// An APIError represents an error that can be returned by the API
type APIError struct {
	// The error code
	Code int
	// The error description
	Description string
}

// An APIMethod represents an API method
type APIMethod struct {
	// The HTTP verb to use for this method (e.g. "GET" or "POST")
	Verb string `json:"method"`
	// A list of the stuff we need to pass to the method
	Input []string
	// A list of what is returned. The API is not consistent here so we don't
	// implement this.
	// https://groups.google.com/forum/#!topic/pcomp15/qSJAm0924Ko
	//Output []string
	// The possible errors
	Errors []APIError
	// An human-readable description of this method
	Description string
}

// APIInfo represents info about the API, returned by… the API itself
type APIInfo struct {
	// All the public methods
	Doc map[string]APIMethod
}
