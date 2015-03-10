package api

import (
	"errors"
	"fmt"
	"github.com/franela/goblin"
	"github.com/franela/goreq"
	o "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"testing"
)

func NewFakeHTTPSJSONServer() *httptest.Server {
	return httptest.NewTLSServer(http.HandlerFunc(func(
		w http.ResponseWriter, r *http.Request) {

		method := r.Method
		path := r.URL.Path

		switch path {
		case "/0/empty":
			w.WriteHeader(200)

		case "/0/empty-object":
			w.WriteHeader(200)
			fmt.Fprint(w, "{}")

		case "/0/method":
			w.WriteHeader(200)
			fmt.Fprintf(w, `{"method": "%s"}`, method)
		}
	}))
}

func TestIOResponses(t *testing.T) {

	g := goblin.Goblin(t)

	o.RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	errDummy := errors.New("a dummy error")

	g.Describe("simpleResponse", func() {
		g.Describe(".IsError()", func() {
			g.It("Should return true if Status is \"error\"", func() {
				sr := simpleResponse{Status: "error"}
				o.Expect(sr.IsError()).To(o.BeTrue())
			})

			g.It("Should return false if Status is \"completed\"", func() {
				sr := simpleResponse{Status: "completed"}
				o.Expect(sr.IsError()).To(o.BeFalse())
			})
		})

		g.Describe(".Error()", func() {
			g.It("Should return nil if .IsError() is false", func() {
				sr := simpleResponse{Status: "completed"}
				o.Expect(sr.IsError()).To(o.BeFalse())
				o.Expect(sr.Error()).To(o.BeNil())
			})

			g.It("Should return an error from the code in .Response", func() {
				sr := simpleResponse{Status: "error"}
				sr.Response.Error_code = 61760457
				o.Expect(sr.IsError()).To(o.BeTrue())
				o.Expect(sr.Error()).To(o.Equal(ErrWrongCmd))
			})
		})
	})

	g.Describe("Body", func() {

		var ts *httptest.Server
		var h *Httclient

		g.Before(func() { ts = NewFakeHTTPSJSONServer() })
		g.After(func() { ts.Close() })

		g.BeforeEach(func() {
			h = NewHTTClient()
			h.baseURL = ts.URL
		})

		g.Describe(".IsEmpty()", func() {
			g.It("Should return true if .Content is nil", func() {
				b := Body{Content: nil}
				o.Expect(b.IsEmpty()).To(o.BeTrue())
			})

			g.It("Should return false if .Content is not nil", func() {
				b := Body{Content: &goreq.Body{}}
				o.Expect(b.IsEmpty()).To(o.BeFalse())
			})
		})

		g.Describe(".Error()", func() {
			g.It("Should return ErrEmptyBody if .IsEmpty()", func() {
				b := Body{Content: nil}
				o.Expect(b.IsEmpty()).To(o.BeTrue())
				o.Expect(b.Error()).To(o.Equal(ErrEmptyBody))
			})

			g.It("Should return .err if it's not nil and the body is empty", func() {
				b := Body{Content: nil, err: errDummy}
				o.Expect(b.IsEmpty()).To(o.BeTrue())
				o.Expect(b.Error()).To(o.Equal(errDummy))
			})

			g.It("Should return .err if it's not nil and the body is not empty", func() {
				b := Body{Content: &goreq.Body{}, err: errDummy}
				o.Expect(b.IsEmpty()).To(o.BeFalse())
				o.Expect(b.Error()).To(o.Equal(errDummy))
			})

			g.It("Should return nil if there's a body and no error", func() {
				b := Body{Content: &goreq.Body{}}
				o.Expect(b.IsEmpty()).To(o.BeFalse())
				o.Expect(b.Error()).To(o.BeNil())
			})
		})

		g.Describe(".FromJSONTo(...)", func() {
			g.It("Shouldn't return nil if .Error() is true", func() {
				target := struct{ foo int }{}
				b := Body{Content: nil, err: errDummy}
				o.Expect(b.FromJSONTo(&target)).To(o.Equal(errDummy))
			})

			g.It("Should return ErrEmptyBody if the body is nil", func() {
				target := struct{ foo int }{}
				b := Body{Content: nil}
				o.Expect(b.FromJSONTo(&target)).To(o.Equal(ErrEmptyBody))
			})

			g.It("Should return an error if the body is empty", func() {
				b := h.call(get, "/empty", nil)
				target := struct{ foo int }{}
				o.Expect(b.FromJSONTo(&target)).NotTo(o.BeNil())
			})

			g.It("Should return nil if the body hasn't any field", func() {
				b := h.call(get, "/empty-object", nil)
				target := struct{ foo int }{foo: 42}
				o.Expect(b.FromJSONTo(&target)).To(o.BeNil())
				o.Expect(target.foo).To(o.Equal(42))
			})

			g.It("Should populate the target struct", func() {
				b := h.call(get, "/method", nil)
				target := struct{ Method string }{}
				o.Expect(b.FromJSONTo(&target)).To(o.BeNil())
				o.Expect(target.Method).To(o.Equal("GET"))
			})
		})

		g.Describe(".Close()", func() {
			g.It("Should return nil if the body is nil", func() {
				b := Body{Content: nil}
				o.Expect(b.Close()).To(o.BeNil())
			})

			g.It("Should return nil if the body was successfully closed", func() {
				b := h.call(get, "/empty-object", nil)
				o.Expect(b.Close()).To(o.BeNil())
			})
		})
	})
}
