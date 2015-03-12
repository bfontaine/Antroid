package api

import (
	"encoding/json"
)

/*
   This file defines structures used to unmarshal JSON responses from the API.
*/

type baseResponse struct {
	Status   string
	Response json.RawMessage
}

type errorResponse struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_msg"`
}

func (e errorResponse) Error() error {
	return errorForCode(e.Code)
}

// Body is a body response from the API
type Body struct {
	Content    *json.RawMessage
	StatusCode int
	err        error
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

// DumpTo takes a pointer to a struct and dumps the body in it
func (b Body) DumpTo(data interface{}) error {
	if err := b.Error(); err != nil {
		return err
	}

	return json.Unmarshal(*b.Content, data)
}

// JSONString returns a JSON string for this body
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

type gameStatusResponse struct {
	Creator        string
	CreationDate   string `json:"creation_date"`
	Teaser         string
	Visibility     string
	NbAntPerPlayer int `json:"nb_ant_per_player"`
	Pace           int
	InitialEnergy  int `json:"initial_energy"`
	InitialAcid    int `json:"initial_acid"`
	Players        []string
	Score          map[string]int
	Status         struct{ Status string }
	Turn           int
}

type playResponse struct {
	Turn         int
	Observations [][]json.RawMessage
}

type visibleAntResponse struct {
	X, Y, Dx, Dy int
	Brain        string
}

type antResponse struct {
	visibleAntResponse
	ID, Energy, Acid int
}

type cellResponse struct {
	X, Y    int
	Content struct {
		Kind  string
		Level string
	}
}

type mapResponse []cellResponse

func (p playResponse) getTurn() (t *Turn, err error) {
	var antResp antResponse
	var mapResp []cellResponse
	var visibleAntsResponse []visibleAntResponse

	t = &Turn{Number: p.Turn, AntsStatuses: []AntStatus{}}

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

		pmap := PartialMap{Cells: make(map[Position]*Cell)}

		for _, cell := range mapResp {
			p := Position{X: cell.X, Y: cell.Y}

			content := cell.Content.Kind
			if cell.Content.Level != "" {
				// the food is represented as a content "food" with a level:
				// "meat", "sugar", etc.
				content = cell.Content.Level
			}

			pmap.Cells[p] = &Cell{
				Pos:     p,
				Content: content,
			}
		}

		ants := []BasicAntStatus{}

		for _, a := range visibleAntsResponse {
			ants = append(ants, BasicAntStatus{
				Pos:   Position{X: a.X, Y: a.Y},
				Dir:   Direction{X: a.Dx, Y: a.Dy},
				Brain: a.Brain,
			})
		}

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

		t.AntsStatuses = append(t.AntsStatuses, ant)
	}

	return
}
