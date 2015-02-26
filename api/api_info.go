package api

type ApiMethod struct {
	Name        string
	Input       []string
	Output      []string
	Errors      []error
	Description string
}

type ApiInfo struct {
	Methods []ApiMethod
}
