package api

// Position is a map position
type Position struct {
	X int
	Y int
}

// A Direction is just a position with offsets instead of absolute variables
type Direction Position

// Cell is a positioned cell
type Cell struct {
	Pos     Position
	Content string
}

type MapInterface interface {
	Width() int
	Height() int
	Cell(int, int) *Cell
}

// PartialMap is a part of a map
type PartialMap struct {
	// we might want to implement this with a map instead
	Cells map[Position]*Cell
}

type Map struct {
	PartialMap

	width, height int
}

func (pm PartialMap) Width() int {
	maxX := -1

	for p, _ := range pm.Cells {
		if p.X > maxX {
			maxX = p.X
		}
	}

	return maxX + 1
}

func (pm PartialMap) Height() int {
	maxY := -1

	for p, _ := range pm.Cells {
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	return maxY + 1
}

func (pm PartialMap) Cell(x, y int) *Cell {
	return pm.Cells[Position{X: x, Y: y}]
}

func (m Map) Width() int  { return m.width }
func (m Map) Height() int { return m.height }
