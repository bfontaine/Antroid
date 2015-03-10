package api

// APIError represents an error that can be returned by the API
type APIError struct {
	// The error code
	Code int
	// The error description
	Description string
}

// APIMethod represents an API method
type APIMethod struct {
	// The HTTP verb to use for this method
	Verb string `json:"method"`
	// A list of the stuff we need to pass to the method
	Input []string
	// A list of what is returned
	// The API is broken here, see
	// https://groups.google.com/forum/#!topic/pcomp15/qSJAm0924Ko
	//Output []string
	// The possible errors
	Errors []APIError
	// An human-readable description of this method
	Description string
}

// APIInfo represents some infos returned by the API
type APIInfo struct {
	// All the public methods
	Doc map[string]APIMethod
}
