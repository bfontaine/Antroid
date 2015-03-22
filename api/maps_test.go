package api

import (
	"github.com/franela/goblin"
	o "github.com/onsi/gomega"
	"testing"
)

func TestMaps(t *testing.T) {

	g := goblin.Goblin(t)

	o.RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("NewPartialMap", func() {
		g.It("Should not return nil", func() {
			o.Expect(NewPartialMap()).NotTo(o.BeNil())
		})

		g.It("Should not return a nil cells array", func() {
			o.Expect(NewPartialMap().Cells).NotTo(o.BeNil())
		})
	})

	g.Describe("PartialMap", func() {
		var pm *PartialMap

		g.BeforeEach(func() {
			pm = NewPartialMap()
		})

		g.Describe(".Width()", func() {
			g.It("Should return 0 on an empty map", func() {
				o.Expect(pm.Width()).To(o.Equal(0))
			})

			g.It("Should return 1 with a cell in (0,0)", func() {
				pos := Position{X: 0, Y: 0}
				pm.Cells[pos] = &Cell{Pos: pos}
				o.Expect(pm.Width()).To(o.Equal(1))
			})

			g.It("Should return 1 with a cell in (0,200)", func() {
				pos := Position{X: 0, Y: 200}
				pm.Cells[pos] = &Cell{Pos: pos}
				o.Expect(pm.Width()).To(o.Equal(1))
			})

			g.It("Should return 201 with a cell in (200,0)", func() {
				pos := Position{X: 200, Y: 0}
				pm.Cells[pos] = &Cell{Pos: pos}
				o.Expect(pm.Width()).To(o.Equal(201))
			})
		})

		g.Describe(".Height()", func() {
			g.It("Should return 0 on an empty map", func() {
				o.Expect(pm.Height()).To(o.Equal(0))
			})

			g.It("Should return 1 with a cell in (0,0)", func() {
				pos := Position{X: 0, Y: 0}
				pm.Cells[pos] = &Cell{Pos: pos}
				o.Expect(pm.Height()).To(o.Equal(1))
			})

			g.It("Should return 201 with a cell in (0,200)", func() {
				pos := Position{X: 0, Y: 200}
				pm.Cells[pos] = &Cell{Pos: pos}
				o.Expect(pm.Height()).To(o.Equal(201))
			})

			g.It("Should return 1 with a cell in (200,0)", func() {
				pos := Position{X: 200, Y: 0}
				pm.Cells[pos] = &Cell{Pos: pos}
				o.Expect(pm.Height()).To(o.Equal(1))
			})
		})

		g.Describe(".Cell(x,y)", func() {
			g.It("Should return nil on an empty map", func() {
				o.Expect(pm.Cell(0, 0)).To(o.BeNil())
			})
		})
	})

	g.Describe("Map", func() {
		var m *Map

		g.BeforeEach(func() {
			m = &Map{
				PartialMap: *NewPartialMap(),
			}
		})

		g.Describe(".Width()", func() {
			g.It("Should return .width", func() {
				m.width = 46
				o.Expect(m.Width()).To(o.Equal(46))
			})
		})

		g.Describe(".Height()", func() {
			g.It("Should return .height", func() {
				m.height = 143
				o.Expect(m.Height()).To(o.Equal(143))
			})
		})

		g.Describe(".Cell(x,y)", func() {
			g.It("Should return nil on an empty map", func() {
				m.width = 0
				m.height = 0
				o.Expect(m.Cell(0, 0)).To(o.BeNil())
			})
		})
	})
}
