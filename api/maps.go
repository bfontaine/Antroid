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
	Pos        Position
	Content    string
	Visibility bool
}

// MapInterface is used for all things that represent maps, i.e. Map (full map)
// and PartialMap
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

// A Map is like a PartialMap but we know its dimensions and all of its content
// This is what is returned by the /log API call
type Map struct {
	PartialMap

	width, height int
}

// NewPartialMap returns a new partial map
func NewPartialMap() *PartialMap {
	return &PartialMap{Cells: make(map[Position]*Cell)}
}

// Width returns the known width of the partial map
func (pm PartialMap) Width() int {
	maxX := -1

	for p := range pm.Cells {
		if p.X > maxX {
			maxX = p.X
		}
	}

	return maxX + 1
}

// Height returns the known height of the partial map
func (pm PartialMap) Height() int {
	maxY := -1

	for p := range pm.Cells {
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	return maxY + 1
}

// Cell returns the cell at (x,y) in the map, or nil if we don't know it
func (pm PartialMap) Cell(x, y int) *Cell {
	return pm.Cells[Position{X: x, Y: y}]
}

// Width returns the width of the map
func (m Map) Width() int { return m.width }

// Height returns the width of the map
func (m Map) Height() int { return m.height }

// Combine modifies the current partial map in-place by adding other partial
// maps to it
func (pm *PartialMap) Combine(maps ...PartialMap) {
	for _, m := range maps {
		for p, c := range m.Cells {
			pm.Cells[p] = c
		}
	}
}

// SetVisibility changes the visibility of all cells in the partial map
func (pm *PartialMap) SetVisibility(v bool) {
	for _, c := range pm.Cells {
		c.Visibility = v
	}
}

// ResetVisibility changes the visibility of all cells in the partial map to
// "false"
func (pm *PartialMap) ResetVisibility() { pm.SetVisibility(false) }

// CombinePartialMaps returns a new PartialMap that is the combination of the
// first one and all the others
func CombinePartialMaps(maps ...PartialMap) *PartialMap {
	pm := NewPartialMap()
	pm.Combine(maps...)
	return pm
}
