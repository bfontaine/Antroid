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
	Verb string `url:"method"`
	// A list of the stuff we need to pass to the method
	Input []string
	// A list of what is returned
	Output []string
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
