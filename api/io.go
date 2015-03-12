package api

import (
	"fmt"
	"github.com/franela/goreq"
	"github.com/google/go-querystring/query"
	"net/http/cookiejar"
	"strings"
)

/* The base URL of all API calls */
const defaultBaseURL = "https://yann.regis-gianas.org/antroid"

/* The API version we support */
const defaultAPIVersion = "0"

/* The User-Agent header we use in all requests */
const defaultUserAgent = "Antroid w/ Go, Cailloux&Fontaine&Galichet&Sagot"

// Httclient is an HTTP client for the API server
type Httclient struct {
	UserAgent string

	baseURL    string
	apiVersion string

	cookies *cookiejar.Jar

	debug bool
}

// NewHTTClient creates a new HTTP client.
func NewHTTClient() *Httclient {
	jar, _ := cookiejar.New(nil)

	return &Httclient{
		UserAgent:  defaultUserAgent,
		baseURL:    defaultBaseURL,
		apiVersion: defaultAPIVersion,
		cookies:    jar,
	}
}

// Return an absolute URL for a given call.
func (h *Httclient) makeAPIURL(call string) string {
	return fmt.Sprintf("%s/%s%s", h.baseURL, h.apiVersion, string(call))
}

// Return the appropriate error for a given HTTP code
func getError(code int) error {
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
func (h *Httclient) call(method, call string, data interface{}) (b *Body) {
	b = &Body{}

	req := goreq.Request{
		Uri:       h.makeAPIURL(call),
		Method:    string(method),
		Accept:    "application/json",
		UserAgent: h.UserAgent,
		// the server uses a self-signed certificate
		Insecure:  true,
		ShowDebug: h.debug,

		CookieJar: h.cookies,
	}

	if method == "GET" {
		// goreq will encode everything for us
		req.QueryString = data
	} else if data != nil {
		// we need to encode our values because the server doesn't accept JSON
		// in requests.
		values, err := query.Values(data)

		if err != nil {
			b.err = err
			return
		}

		// set all keys to lower-case and remove multiple values (e.g. a=2&a=3)
		for k, v := range values {
			if len(k) > 0 && len(v) > 0 && k[0] >= 'A' && k[0] <= 'Z' {
				values.Set(strings.ToLower(k), v[0])
				values.Del(k)
			}
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

	if res.Body != nil {
		defer res.Body.Close()
	}

	if b.StatusCode = res.StatusCode; b.StatusCode != 200 {
		b.err = getError(b.StatusCode)
		return
	}

	var resp baseResponse

	b.err = res.Body.FromJsonTo(&resp)

	if h.debug {
		fmt.Printf("\n==> %s\n\n", resp)
	}

	if b.err != nil {
		return
	}

	b.Content = &resp.Response

	switch resp.Status {
	case "completed":
		// nothing here, everything's fine
	case "error":
		var errorResp errorResponse

		if b.err = b.DumpTo(&errorResp); b.err == nil {
			b.err = errorResp.Error()
		}

	default:
		b.err = ErrUnknown
	}

	return
}

const (
	get  = "GET"
	post = "POST"
)

func (h *Httclient) get(call string, data interface{}) *Body {
	return h.call(get, call, data)
}

func (h *Httclient) post(call string, data interface{}) *Body {
	return h.call(post, call, data)
}

/*
   Each method below perform a call to one endpoint. We expose them instead of
   the generic .call method to be able to type-check the parameters of each
   call.
*/

// CallAPI performs a call to /api.
func (h *Httclient) CallAPI() *Body { return h.get("/api", NoParams{}) }

// CallAuth performs a call to /auth.
func (h *Httclient) CallAuth(params UserCredentialsParams) *Body {
	return h.post("/auth", params)
}

// CallCreate performs a call to /create.
func (h *Httclient) CallCreate(params GameSpecParams) *Body {
	return h.get("/create", params)
}

// CallDestroy performs a call to /destroy.
func (h *Httclient) CallDestroy(params GameIDParams) *Body {
	return h.get("/destroy", params)
}

// CallGames performs a call to /games.
func (h *Httclient) CallGames() *Body {
	return h.get("/games", NoParams{})
}

// CallJoin performs a call to /join.
func (h *Httclient) CallJoin(params GameIDParams) *Body {
	return h.get("/join", params)
}

// CallLog performs a call to /log.
func (h *Httclient) CallLog(params GameIDParams) *Body {
	return h.get("/log", params)
}

// CallLogout performs a call to /logout.
func (h *Httclient) CallLogout() *Body {
	return h.get("/logout", NoParams{})
}

// CallPlay performs a call to /play.
func (h *Httclient) CallPlay(params PlayParams) *Body {
	return h.get("/play", params)
}

// CallRegister performs a call to /register.
func (h *Httclient) CallRegister(params UserCredentialsParams) *Body {
	return h.post("/register", params)
}

// CallShutdown performs a call to /shutdown.
func (h *Httclient) CallShutdown(params GenericIDParams) *Body {
	return h.get("/shutdown", params)
}

// CallStatus performs a call to /status.
func (h *Httclient) CallStatus(params GameIDParams) *Body {
	return h.get("/status", params)
}

// CallWhoAmI performs a call to /whoami.
func (h *Httclient) CallWhoAmI() *Body {
	return h.get("/whoami", NoParams{})
}
