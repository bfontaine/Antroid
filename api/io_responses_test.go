package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/franela/goblin"
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
			fmt.Fprint(w, `{"status":"completed", "response":{}}`)

		case "/0/method":
			w.WriteHeader(200)
			fmt.Fprintf(w, `
                {"status":"completed", "response":{"method": "%s"}}`, method)
		}
	}))
}

func TestIOResponses(t *testing.T) {

	g := goblin.Goblin(t)

	o.RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	errDummy := errors.New("a dummy error")

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
				b := Body{Content: &json.RawMessage{}}
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
				b := Body{Content: &json.RawMessage{}, err: errDummy}
				o.Expect(b.IsEmpty()).To(o.BeFalse())
				o.Expect(b.Error()).To(o.Equal(errDummy))
			})

			g.It("Should return nil if there's a body and no error", func() {
				b := Body{Content: &json.RawMessage{}}
				o.Expect(b.IsEmpty()).To(o.BeFalse())
				o.Expect(b.Error()).To(o.BeNil())
			})
		})

		g.Describe(".DumpTo(s)", func() {
			g.It("Shouldn't return nil if .Error() is true", func() {
				target := struct{ foo int }{}
				b := Body{Content: nil, err: errDummy}
				o.Expect(b.DumpTo(&target)).To(o.Equal(errDummy))
			})

			g.It("Should return ErrEmptyBody if the body is nil", func() {
				target := struct{ foo int }{}
				b := Body{Content: nil}
				o.Expect(b.DumpTo(&target)).To(o.Equal(ErrEmptyBody))
			})

			g.It("Should return an error if the body is empty", func() {
				b := h.call(get, "/empty", nil)
				target := struct{ foo int }{}
				o.Expect(b.DumpTo(&target)).NotTo(o.BeNil())
			})

			g.It("Should return nil if the body hasn't any field", func() {
				b := h.call(get, "/empty-object", nil)
				target := struct{ foo int }{foo: 42}
				o.Expect(b.DumpTo(&target)).To(o.BeNil())
				o.Expect(target.foo).To(o.Equal(42))
			})

			g.It("Should populate the target struct", func() {
				b := h.call(get, "/method", nil)
				target := struct{ Method string }{}
				o.Expect(b.DumpTo(&target)).To(o.BeNil())
				o.Expect(target.Method).To(o.Equal("GET"))
			})
		})
	})
}
