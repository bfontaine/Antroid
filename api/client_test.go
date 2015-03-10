package api

import (
	"fmt"
	"github.com/franela/goblin"
	o "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"testing"
)

func NewFakeAPIServer() *httptest.Server {
	return httptest.NewTLSServer(http.HandlerFunc(func(
		w http.ResponseWriter, r *http.Request) {

		route := fmt.Sprintf("%s %s", r.Method, r.URL.Path)

		r.ParseForm()

		switch route {
		case "GET /0/api":
			w.WriteHeader(200)
			fmt.Fprint(w, `{
              "status": "completed",
              "response": {
              "doc": {
                "m1": {
                  "method": "post",
                  "input": [ "i1 : string", "i2 : string" ],
                  "output": [],
                  "errors": [
                    { "code": 332299703, "description": "USER_ALREADY_EXISTS" },
                    { "code": 621433138, "description": "INVALID_LOGIN" }
                  ],
                  "description": "Do something"
                },
                "m2": {
                  "method": "get",
                  "input": [],
                  "output": [ "foo : string" ],
                  "errors": [],
                  "description": "Do something else"
                }
              }
            }
          }`)
		}

		// TODO mock the API
	}))
}

func TestClient(t *testing.T) {

	g := goblin.Goblin(t)

	o.RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("NewClient", func() {
		g.It("Should not return nil", func() {
			o.Expect(NewClient()).NotTo(o.BeNil())
		})
	})

	g.Describe("Client", func() {
		var ts *httptest.Server
		var c *Client

		g.BeforeEach(func() {
			ts = NewFakeAPIServer()
			c, _ = NewClient()
			c.http.baseURL = ts.URL
		})

		g.AfterEach(func() { ts.Close() })

		g.Describe(".getUserCredentialsParams", func() {
			g.It("Should initially return an empty UserCredentialsParams", func() {
				p := UserCredentialsParams{}
				o.Expect(c.getUserCredentialsParams()).To(o.Equal(p))
			})

			g.It("Should return an UserCredentialsParams w/ username/password", func() {
				c.username = "foo"
				c.password = "secret"
				p := UserCredentialsParams{Login: "foo", Password: "secret"}
				o.Expect(c.getUserCredentialsParams()).To(o.Equal(p))
			})
		})

		g.Describe(".Authenticated", func() {
			g.It("Should be initially false", func() {
				o.Expect(c.Authenticated()).To(o.BeFalse())
			})
		})

		g.Describe(".APIInfo", func() {
			g.It("Should return API infos", func() {
				info, err := c.APIInfo()

				o.Expect(err).To(o.BeNil())
				o.Expect(info).NotTo(o.BeNil())
				o.Expect(info.Doc).NotTo(o.BeNil())
			})

			g.It("Should populate .Doc with all methods", func() {
				info, err := c.APIInfo()

				o.Expect(err).To(o.BeNil())
				o.Expect(info).NotTo(o.BeNil())
				o.Expect(info.Doc).NotTo(o.BeNil())
				o.Expect(len(info.Doc)).To(o.Equal(2))

				m1, m2 := info.Doc["m1"], info.Doc["m2"]

				o.Expect(m1).To(o.Equal(APIMethod{
					Verb:   "post",
					Input:  []string{"i1 : string", "i2 : string"},
					Output: []string{},
					Errors: []APIError{
						APIError{Code: 332299703,
							Description: "USER_ALREADY_EXISTS"},
						APIError{Code: 621433138,
							Description: "INVALID_LOGIN"},
					},
					Description: "Do something",
				}))

				o.Expect(m2).To(o.Equal(APIMethod{
					Verb:        "get",
					Input:       []string{},
					Output:      []string{"foo : string"},
					Errors:      []APIError{},
					Description: "Do something else",
				}))
			})

			// TODO
		})
	})
}
