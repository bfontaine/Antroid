package api

import (
	"fmt"
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

	g.Describe("Command#String()", func() {
		g.It("Should return {.Ant}:{.Cmd}", func() {
			o.Expect(Command{Ant: 1, Cmd: "foo"}.String()).To(o.Equal("1:foo"))
		})
	})

	g.Describe("Commands#String()", func() {
		g.It("Should return an empty string if there're no commands", func() {
			o.Expect(Commands{}.String()).To(o.Equal(""))
		})

		g.It("Should return a command if there's only one", func() {
			cmd := Command{Ant: 42, Cmd: "qux"}
			o.Expect(Commands{cmd}.String()).To(o.Equal(cmd.String()))
		})

		g.It("Should return a CSV line if there're multiple commands", func() {
			cmd1 := Command{Ant: 42, Cmd: "qux"}
			cmd2 := Command{Ant: 17, Cmd: "bar"}
			expected := fmt.Sprintf("%s,%s,%s", cmd1, cmd2, cmd1)
			o.Expect(Commands{cmd1, cmd2, cmd1}.String()).To(o.Equal(expected))
		})
	})
}
