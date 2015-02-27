package api

/*
   This file defines structures used to unmarshal JSON responses from the API.
*/

type registerResponse struct {
	Status   string
	Response struct {
		Error_code int
		Error_msg  string
	}
}

type whoAmIResponse struct {
	Status   string
	Response struct{ Status string }
}
