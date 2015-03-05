package api

import (
	"github.com/franela/goblin"
	o "github.com/onsi/gomega"
	"testing"
)

func TestIOResponses(t *testing.T) {

	g := goblin.Goblin(t)

	o.RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("simpleResponse", func() {
		g.Describe(".IsError()", func() {
			g.It("Should return true if Status is \"error\"", func() {
				sr := simpleResponse{Status: "error"}

				o.Expect(sr.IsError()).To(o.BeTrue())
			})
		})
	})
}
