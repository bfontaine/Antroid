package api

import (
	"github.com/franela/goblin"
	o "github.com/onsi/gomega"
	"testing"
)

func TestGameSpec(t *testing.T) {

	g := goblin.Goblin(t)

	o.RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("IntRange", func() {
		g.It("Shouldn't return nil", func() {
			o.Expect(IntRange(0, 2)).NotTo(o.BeNil())
		})
	})

	g.Describe("intRange", func() {
		g.Describe(".Include(n)", func() {
			var r intRange

			g.BeforeEach(func() { r = IntRange(1, 42) })

			g.It("Should return false if n < .min", func() {
				o.Expect(r.Include(-42)).To(o.BeFalse())
			})

			g.It("Should return false if n > .max", func() {
				o.Expect(r.Include(101)).To(o.BeFalse())
			})

			g.It("Should return true if n = .min", func() {
				o.Expect(r.Include(1)).To(o.BeTrue())
			})

			g.It("Should return true if n = .max", func() {
				o.Expect(r.Include(42)).To(o.BeTrue())
			})

			g.It("Should return true if .min < n < .max", func() {
				o.Expect(r.Include(30)).To(o.BeTrue())
			})
		})
	})

	g.Describe("GameSpec", func() {
		g.Describe(".toParams()", func() {
			g.It("Should set .Users to \"all\" if the game is public", func() {
				gs := GameSpec{Public: true}
				o.Expect(gs.toParams().Users).To(o.Equal("all"))
			})

			g.It("Should set .Users to a comma-separated list of players "+
				"if the game is private", func() {

				gs := GameSpec{Players: []string{"a", "b"}}
				o.Expect(gs.toParams().Users).To(o.Equal("a,b"))
			})
		})

		g.Describe(".Validate()", func() {
			g.It("Should return false if the game is private without players", func() {
				gs := GameSpec{}
				o.Expect(gs.Validate()).To(o.BeFalse())
			})
			g.It("Should return false if the game is public with players", func() {
				gs := GameSpec{Public: true, Players: []string{"a"}}
				o.Expect(gs.Validate()).To(o.BeFalse())
			})
			g.It("Should return false if the pace is negative", func() {
				gs := GameSpec{Public: true, Pace: -12}
				o.Expect(gs.Validate()).To(o.BeFalse())
			})
			g.It("Should return true all parameters are in the range", func() {
				gs := GameSpec{
					Public:        true,
					Pace:          2,
					Turns:         2,
					AntsPerPlayer: 2,
					MaxPlayers:    2,
					MinPlayers:    2,
					InitialEnergy: 40,
					InitialAcid:   40,
				}
				o.Expect(gs.Validate()).To(o.BeTrue())
			})
		})
	})
}
