package api

// This file describes structures used to unparse JSON responses from the API.
// This is the opposite of `api/params.go`.

import "encoding/json"

// A baseResponse is the mother of all responses. It's used by the low-level
// HTTP(S) client to get the server's status. The .Response field is then
// passed to the high-level client which unparses it in the right structs.
type baseResponse struct {
	Status   string
	Response json.RawMessage
}

// An error response has a code and a message
type errorResponse struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_msg"`
}

// Return a Go error from the code of an error
func (e errorResponse) Error() error {
	return errorForCode(e.Code)
}

// Body is a body response from the API
type Body struct {
	// its content is a raw JSON string
	Content *json.RawMessage
	// status code from the server
	StatusCode int
	// any error we got when performing the request
	err error
}

// IsEmpty returns true if the body response is empty
func (b Body) IsEmpty() bool {
	return b.Content == nil
}

// Error returns any error with this body
func (b Body) Error() (err error) {
	if b.err != nil {
		return b.err
	}

	if b.IsEmpty() {
		err = ErrEmptyBody
	}

	return
}

// DumpTo takes a pointer to a struct and dumps the body in it. It'll return an
// error if it can't dump or if the body's .Error() is non-nil.
func (b Body) DumpTo(data interface{}) error {
	if err := b.Error(); err != nil {
		return err
	}

	return json.Unmarshal(*b.Content, data)
}

// JSONString returns a JSON string for this body. This is used for debugging
// purposes.
func (b Body) JSONString() string {
	return string(*b.Content)
}

// ensureEmptyResponse tries to parse the body as an empty JSON object. It'll
// return an error if .Error() is non-nil
func (b Body) ensureEmptyResponse() (err error) {
	if err = b.Error(); err == nil {
		err = b.DumpTo(&struct{}{})
	}

	return
}

// a gameStatusResponse represents the response from the server when querying
// for the status of a game
type gameStatusResponse struct {
	Creator        string
	CreationDate   string `json:"creation_date"`
	Teaser         string
	Visibility     json.RawMessage
	NbAntPerPlayer int `json:"nb_ant_per_player"`
	Pace           int
	InitialEnergy  int `json:"initial_energy"`
	InitialAcid    int `json:"initial_acid"`
	Players        []string
	Score          map[string]int
	Status         struct{ Status string }
	Turn           int
}

// a playResponse is a partially parsed result from an API call to /play
type playResponse struct {
	Turn         int
	Observations [][]json.RawMessage
}

// visibleAntResponse is a part of a response which describes a visible ant
type visibleAntResponse struct {
	X, Y, Dx, Dy int
	Brain        string
}

// visibleAntResponse is a part of a response which describes one of our ants
type antResponse struct {
	visibleAntResponse
	ID, Energy, Acid int
}

// cellResponse is a part of a response which describes a map cell
type cellResponse struct {
	X, Y    int
	Content struct {
		Kind  string
		Level string
	}
}

// A map response is a list of cellResponses
type mapResponse []cellResponse

// getTurn returns a Turn object extracted from a playResponse
func (p playResponse) getTurn() (t *Turn, err error) {
	var antResp antResponse
	var mapResp []cellResponse
	var visibleAntsResponse []visibleAntResponse

	// create an empty turn
	t = &Turn{Number: p.Turn, AntsStatuses: []AntStatus{}}

	// for each one of our ants
	for _, obs := range p.Observations {

		// ant info
		if err = json.Unmarshal(obs[0], &antResp); err != nil {
			return
		}

		// map info
		if err = json.Unmarshal(obs[1], &mapResp); err != nil {
			return
		}

		// visible ants info
		if err = json.Unmarshal(obs[2], &visibleAntsResponse); err != nil {
			return
		}

		// create the partial map for this turn
		pmap := NewPartialMap()

		// populate it with all cells
		for _, cell := range mapResp {
			p := Position{X: cell.X, Y: cell.Y}

			content := cell.Content.Kind
			if cell.Content.Level != "" {
				// the food is represented as a content "food" with a level:
				// "meat", "sugar", etc.
				content = cell.Content.Level
			}

			// add a cell at its position
			pmap.Cells[p] = &Cell{
				Pos:     p,
				Content: content,
			}
		}

		// populate the visible ants list
		ants := []BasicAntStatus{}

		for _, a := range visibleAntsResponse {
			ants = append(ants, BasicAntStatus{
				Pos:   Position{X: a.X, Y: a.Y},
				Dir:   Direction{X: a.Dx, Y: a.Dy},
				Brain: a.Brain,
			})
		}

		// create this ant's status struct
		ant := AntStatus{
			BasicAntStatus: BasicAntStatus{
				Pos:   Position{X: antResp.X, Y: antResp.Y},
				Dir:   Direction{X: antResp.Dx, Y: antResp.Dy},
				Brain: antResp.Brain,
			},

			ID:          antResp.ID,
			Energy:      antResp.Energy,
			Acid:        antResp.Acid,
			Vision:      pmap,
			VisibleAnts: ants,
		}

		// add the ant status to our list
		t.AntsStatuses = append(t.AntsStatuses, ant)
	}

	return
}
