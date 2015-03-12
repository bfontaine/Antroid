package api

// Commands is a comma-separated list of commands we give to an ant
type Commands string

// Combine returns a new PartialMap that is the combination of the first one
// and all the others
func (pm PartialMap) Combine(maps ...PartialMap) PartialMap {
	for _, m := range maps {
		for p, c := range m.Cells {
			pm.Cells[p] = c
		}
	}

	return pm
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
