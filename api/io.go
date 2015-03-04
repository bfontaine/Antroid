package api

import (
	"fmt"
	"github.com/franela/goreq"
	"github.com/google/go-querystring/query"
	"net/http/cookiejar"
)

/* The base URL of all API calls */
const BASE_URL = "https://yann.regis-gianas.org/antroid"

/* The API version we support */
const API_VERSION = "0"

/* The User-Agent header we use in all requests */
const USER_AGENT = "Antroid w/ Go, Cailloux&Fontaine&Galichet&Sagot"

// An HTTP client for the API server
type httclient struct {
	UserAgent string

	baseUrl    string
	apiVersion string

	cookies *cookiejar.Jar
}

// Create a new HTTP client.
func NewHTTClient() httclient {
	jar, _ := cookiejar.New(nil)

	return httclient{
		UserAgent:  USER_AGENT,
		baseUrl:    BASE_URL,
		apiVersion: API_VERSION,
		cookies:    jar,
	}
}

// Return an absolute URL for a given call.
func (h *httclient) makeApiUrl(call string) string {
	return fmt.Sprintf("%s/%s%s", h.baseUrl, h.apiVersion, string(call))
}

// Return the appropriate error for a given HTTP code
func (h *httclient) getError(code int) error {
	switch code / 100 {
	case 4:
		return Err4XX
	case 5:
		return Err5XX
	}

	return nil
}

// Make an HTTP call to the remote server and return its response body.
// Don't forget to close it if it's not nil.
func (h *httclient) call(method, call string, data interface{}) (b *Body) {
	req := goreq.Request{
		Uri:       h.makeApiUrl(call),
		Method:    string(method),
		Accept:    "application/json",
		UserAgent: h.UserAgent,
		// the server uses a self-signed certificate
		Insecure: true,
		//ShowDebug: true,

		CookieJar: h.cookies,
	}

	if method == "GET" {
		// goreq will encode everything for us
		req.QueryString = data
	} else {
		// we need to encode our values because the server doesn't accept JSON in
		// requests.
		values, err := query.Values(data)

		if err != nil {
			b.err = err
			return
		}

		queryString := values.Encode()
		req.ContentType = "application/x-www-form-urlencoded"
		req.Body = queryString
	}

	res, err := req.Do()

	if err != nil {
		b.err = err
		return
	}

	b.Content = res.Body

	if res.StatusCode != 200 {
		err = h.getError(res.StatusCode)
	}

	return
}

const (
	get  = "GET"
	post = "POST"
)

/*
   Each method below perform a call to one endpoint. We expose them instead of
   the generic .call method to be able to type-check the parameters of each
   call.
*/

// Perform a call to /api.
func (h *httclient) CallApi() *Body {
	return h.call(get, "/api", NoParams{})
}

// Perform a call to /auth.
func (h *httclient) CallAuth(params UserCredentialsParams) *Body {
	return h.call(post, "/auth", params)
}

// Perform a call to /create.
func (h *httclient) CallCreate(params GameSpecParams) *Body {
	return h.call(get, "/create", params)
}

// Perform a call to /destroy.
func (h *httclient) CallDestroy(params GameIdParams) *Body {
	return h.call(get, "/destroy", params)
}

// Perform a call to /games.
func (h *httclient) CallGames() *Body {
	return h.call(get, "/games", NoParams{})
}

// Perform a call to /join.
func (h *httclient) CallJoin(params GameIdParams) *Body {
	return h.call(get, "/join", params)
}

// Perform a call to /log.
func (h *httclient) CallLog(params GameIdParams) *Body {
	return h.call(get, "/log", params)
}

// Perform a call to /logout.
func (h *httclient) CallLogout() *Body {
	return h.call(get, "/logout", NoParams{})
}

// Perform a call to /play.
func (h *httclient) CallPlay(params PlayParams) *Body {
	return h.call(get, "/play", params)
}

// Perform a call to /register.
func (h *httclient) CallRegister(params UserCredentialsParams) *Body {
	return h.call(post, "/register", params)
}

// Perform a call to /shutdown.
func (h *httclient) CallShutdown(params GenericIdParams) *Body {
	return h.call(get, "/shutdown", params)
}

// Perform a call to /status.
func (h *httclient) CallStatus(params GameIdParams) *Body {
	return h.call(get, "/status", params)
}

// Perform a call to /whoami.
func (h *httclient) CallWhoAmI() *Body {
	return h.call(get, "/whoami", NoParams{})
}
