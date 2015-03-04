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
	return httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		path := r.URL.Path

		r.ParseForm()

		switch path {
		case "/0/method":
			w.WriteHeader(200)
			fmt.Fprint(w, fmt.Sprintf("%v %v", method, path))

		case "/0/idontexist":
			w.WriteHeader(404)
			fmt.Fprint(w, "nope")

		case "/0/geturlparams":
			if method == "GET" {
				w.WriteHeader(200)
				fmt.Fprint(w, fmt.Sprintf("%v %s", method, r.Form.Encode()))
			}

		case "/0/getpostparams":
			if method == "POST" {
				w.WriteHeader(200)
				fmt.Fprint(w, fmt.Sprintf("%v %s", method, r.PostForm.Encode()))
			}

			// TODO mock the real API
		}

	}))
}

func TestIO(t *testing.T) {

	g := goblin.Goblin(t)

	o.RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("NewHTTClient", func() {
		g.It("Should not return nil", func() {
			o.Expect(NewHTTClient()).NotTo(o.BeNil())
		})
	})

	g.Describe("getError", func() {
		g.It("Should return an Err4XX if the code is 4XX", func() {
			o.Expect(getError(400)).To(o.Equal(Err4XX))
			o.Expect(getError(403)).To(o.Equal(Err4XX))
			o.Expect(getError(404)).To(o.Equal(Err4XX))
		})

		g.It("Should return an Err5XX if the code is 5XX", func() {
			o.Expect(getError(500)).To(o.Equal(Err5XX))
		})
	})

	g.Describe("Httclient", func() {
		g.Describe(".call", func() {
			var ts *httptest.Server
			var h *Httclient

			g.Before(func() { ts = NewFakeAPIServer() })
			g.After(func() { ts.Close() })

			g.BeforeEach(func() {
				h = NewHTTClient()
				h.baseURL = ts.URL
			})

			g.It("Should call the remote server", func() {
				b := h.call(get, "/method", nil)

				o.Expect(b).NotTo(o.BeNil())
				o.Expect(b.Error()).To(o.BeNil())
				o.Expect(b.IsEmpty()).To(o.BeFalse())
				o.Expect(b.StatusCode).To(o.Equal(200))
				o.Expect(b.Content.ToString()).To(o.Equal("GET /0/method"))
			})

			g.It("Should use POST if it was given as the method", func() {
				b := h.call(post, "/method", nil)

				o.Expect(b).NotTo(o.BeNil())
				o.Expect(b.Error()).To(o.BeNil())
				o.Expect(b.IsEmpty()).To(o.BeFalse())
				o.Expect(b.StatusCode).To(o.Equal(200))
				o.Expect(b.Content.ToString()).To(o.Equal("POST /0/method"))
			})

			g.It("Should set body.err to Err4XX if the status code is 4XX", func() {
				b := h.call(get, "/idontexist", nil)

				o.Expect(b).NotTo(o.BeNil())
				o.Expect(b.Error()).To(o.Equal(Err4XX))
				o.Expect(b.StatusCode).To(o.Equal(404))
			})

			g.It("Should send parameters in the URL for GET requests", func() {
				b := h.call(get, "/geturlparams", struct {
					Param string
				}{"foo"})

				o.Expect(b).NotTo(o.BeNil())
				o.Expect(b.Error()).To(o.BeNil())
				o.Expect(b.IsEmpty()).To(o.BeFalse())
				o.Expect(b.StatusCode).To(o.Equal(200))
				o.Expect(b.Content.ToString()).To(o.Equal("GET param=foo"))
			})

			g.It("Should send parameters in the body for POST requests", func() {
				b := h.call(post, "/getpostparams", struct {
					Param string
				}{"foo"})

				o.Expect(b).NotTo(o.BeNil())
				o.Expect(b.Error()).To(o.BeNil())
				o.Expect(b.IsEmpty()).To(o.BeFalse())
				o.Expect(b.StatusCode).To(o.Equal(200))
				o.Expect(b.Content.ToString()).To(o.Equal("POST param=foo"))
			})
		})
	})
}
