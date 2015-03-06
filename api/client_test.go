package api

import (
	"github.com/franela/goblin"
	o "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"testing"
)

func NewFakeAPIServer() *httptest.Server {
	return httptest.NewTLSServer(http.HandlerFunc(func(
		w http.ResponseWriter, r *http.Request) {

		//method := r.Method
		//path := r.URL.Path

		r.ParseForm()

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
			ts = NewFakeHTTPSServer()
			c, _ = NewClient()
			c.http.baseURL = ts.URL
		})

		g.AfterEach(func() { ts.Close() })

		g.Describe(".getUserCredentialsParams", func() {
			// TODO
		})

		g.Describe(".Authenticated", func() {
			g.It("Should be initially false", func() {
				o.Expect(c.Authenticated()).To(o.BeFalse())
			})
		})
	})
}
