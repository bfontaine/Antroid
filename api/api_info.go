package api

// An error that can be returned by the API
type ApiError struct {
	// The error code
	Code int
	// The error description
	Description string
}

// An API method
type ApiMethod struct {
	// The HTTP verb to use for this method
	Verb string `url:"method"`
	// A list of the stuff we need to pass to the method
	Input []string
	// A list of what is returned
	Output []string
	// The possible errors
	Errors []ApiError
	// An human-readable description of this method
	Description string
}

// Some info about the API
type ApiInfo struct {
	// All the public methods
	Doc map[string]ApiMethod
}
