package api

// Commands is a comma-separated list of commands we give to an ant
type Commands string

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

	Vision      *PartialMap
	VisibleAnts []BasicAntStatus
}

// OtherVisibleAnts returns an equivalent of AntStatus.VisibleAnts, excluding
// the ant itself. The API server gives us a list *including* the ant (which
// sees itself).
func (as AntStatus) OtherVisibleAnts() []BasicAntStatus {
	var others []BasicAntStatus

	for _, ant := range as.VisibleAnts {
		if ant.Pos == as.Pos && ant.Dir == as.Dir && ant.Brain == as.Brain {
			continue
		}

		// Note: this could be optimized, since we sees ourselves only once
		others = append(others, ant)
	}

	return others
}

// Turn describes all the infos we have about a turn
type Turn struct {
	Number int

	AntsStatuses []AntStatus
}

var EmptyTurn = Turn{
	Number:       0,
	AntsStatuses: []AntStatus{},
}
