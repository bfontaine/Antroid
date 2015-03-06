package api

import (
	"github.com/franela/goblin"
	o "github.com/onsi/gomega"
	"testing"
)

func TestErrors(t *testing.T) {

	g := goblin.Goblin(t)

	o.RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("errorForCode", func() {
		g.It("Should return ErrUnknownCode if the error code doesn't exist", func() {
			o.Expect(errorForCode(-42)).To(o.Equal(ErrUnknownCode))
		})
	})
}
