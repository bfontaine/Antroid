package api

type ApiError struct {
	Code        int
	Description string
}

type ApiMethod struct {
	Method      string
	Input       []string
	Output      []string
	Errors      []ApiError
	Description string
}

type ApiInfo struct {
	Doc map[string]ApiMethod
}
