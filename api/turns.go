package api

// Commands is a comma-separated list of commands we give to an ant
type Commands string

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

// PartialMap is a part of a map
type PartialMap struct {
	Cells []Cell
}

// BasicAntStatus describes an ant from which we don't know much info
type BasicAntStatus struct {
	Pos   Position
	Dir   Direction
	Brain string
}

// AntStatus describes an ant
type AntStatus struct {
	BasicAntStatus

	ID     int
	Energy int
	Acid   int

	Vision      PartialMap
	VisibleAnts []BasicAntStatus
}

// Turn describes all the infos we have about a turn
type Turn struct {
	Number int

	AntsStatuses []AntStatus
}
