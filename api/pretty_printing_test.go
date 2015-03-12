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

	g.Describe("Commands#String()", func() {
		g.It("Should return an empty string if there're no commands", func() {
			o.Expect(Commands("").String()).To(o.Equal(""))
		})

		g.It("Should return a command if there's only one", func() {
			o.Expect(Commands("1:foo").String()).To(o.Equal("1:foo"))
		})
	})
}
