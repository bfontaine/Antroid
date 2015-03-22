package api

// This file describes what we know after each turn

// Commands is a comma-separated list of commands we give to an ant
type Commands string

// BasicAntStatus describes an ant from which we don't know much info
type BasicAntStatus struct {
	// its position
	Pos Position
	// its direction
	Dir Direction
	// its brain stage (we know "controlled", but they are more states)
	Brain string
}

// Eq tests if this status is equal to another
func (as BasicAntStatus) Eq(other BasicAntStatus) bool {
	return as.Pos.X == other.Pos.X &&
		as.Pos.Y == other.Pos.Y &&
		as.Dir.X == other.Dir.X &&
		as.Dir.Y == other.Dir.Y &&
		as.Brain == other.Brain
}

// AntStatus describes an ant
type AntStatus struct {
	// it inherits all fields from BasicAntStatus
	BasicAntStatus

	// its ID
	ID int
	// its energy level
	Energy int
	// its acid level
	Acid int

	// what she sees around
	Vision *PartialMap
	// all the ants she sees (including itself)
	VisibleAnts []BasicAntStatus
}

// Turn describes all the infos we have about a turn
type Turn struct {
	// its number
	Number int

	// all our ants' statuses
	AntsStatuses []AntStatus
}

// EmptyTurn represents an empty turn with no ants
var EmptyTurn = Turn{
	Number:       0,
	AntsStatuses: []AntStatus{},
}
