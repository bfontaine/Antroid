package api

import (
	"errors"
	"github.com/franela/goblin"
	"github.com/franela/goreq"
	o "github.com/onsi/gomega"
	"testing"
)

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
		})
	})

	g.Describe("Body", func() {
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
	})
}
