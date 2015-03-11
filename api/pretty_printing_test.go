package api

import (
	"github.com/franela/goblin"
	o "github.com/onsi/gomega"
	"testing"
)

func TestPrettyPrinting(t *testing.T) {

	g := goblin.Goblin(t)

	o.RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Game#String()", func() {
		g.It("Should not return an empty string", func() {
			o.Expect(Game{}.String()).NotTo(o.Equal(""))
		})
	})

	g.Describe("GameStatus#String()", func() {
		g.It("Should not return an empty string", func() {
			o.Expect(GameStatus{}.String()).NotTo(o.Equal(""))
		})
	})
}
