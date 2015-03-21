package api

import (
	"bytes"
	"fmt"
)

// A Player represents a local game server connected to the remote one and
// which sends and receive messages from it.
type Player struct {
	Client    *Client
	AIs       *AIPool
	Listeners *ListenersPool

	username, password string

	debug bool

	status     *GameStatus
	turn       *Turn
	partialMap *PartialMap

	done bool
}

// NewPlayer returns a pointer on a new Player
func NewPlayer(username, password string) (p *Player) {
	p = &Player{
		Client:     NewClient(),
		AIs:        NewAIPool(),
		Listeners:  NewListenersPool(),
		username:   username,
		password:   password,
		status:     &GameStatus{},
		turn:       &EmptyTurn,
		partialMap: NewPartialMap(),
	}

	return
}

// SetDebug enables/disables the debug mode
func (p *Player) SetDebug(debug bool) {
	p.Client.SetDebug(debug)
	p.debug = debug
}

// Connect connects the player to the remote server, first trying to register
// its credentials.
func (p *Player) Connect() (err error) {
	// try to register, just in case the credentials don't exist
	err = p.Client.RegisterWithCredentials(p.username, p.password)

	if err != nil && err != ErrUserAlreadyExists {
		return
	}

	return p.Client.Login()
}

// CreateAndJoinGame creates a new game from the given spec and joins it
func (p *Player) CreateAndJoinGame(gs *GameSpec) (err error) {
	var g *Game

	if g, err = p.Client.CreateGame(gs); err != nil {
		return
	}

	err = p.JoinGame(g.Identifier)

	return
}

// JoinGame joins an existing game
func (p *Player) JoinGame(id GameID) (err error) {

	if err = p.Client.JoinGameIdentifier(id); err != nil {
		return
	}

	if p.status, err = p.Client.GetGameIdentifierStatus(id); err != nil {
		return
	}

	firstAnt := true
	var restCmd bytes.Buffer

	// send a "rest" command to all ants for the first turn, just to get all
	// ants' positions
	for i := 0; i < p.status.Game.Spec.AntsPerPlayer; i++ {
		if !firstAnt {
			restCmd.WriteString(",")
		}
		firstAnt = false
		restCmd.WriteString(fmt.Sprintf("%d:rest", i))
	}

	commands := Commands(restCmd.String())
	p.turn, err = p.Client.PlayIdentifier(p.status.Identifier, commands)

	if err != nil {
		return
	}

	p.startPlugins()

	return
}

// PlayTurn sends the game status to all AIs and gets their feedback before
// sending everything to the remote server
func (p *Player) PlayTurn() (done bool, err error) {
	p.sendTurnStatusToPlugins()
	err = p.playTurn()
	done = p.done

	// end of game
	if done && err == ErrGameNotPlaying {
		err = nil
	}

	return
}

// Quit stops all AIs and logout the player from the remote server
func (p *Player) Quit() error {
	p.AIs.Stop()
	p.Listeners.Stop()
	return p.Client.Logout()
}

// PrintScores prints the current scores
func (p *Player) PrintScores() {
	var usernameMaxSize int

	for username := range p.status.Score {
		usernameSize := len(username)

		if usernameSize > usernameMaxSize {
			usernameMaxSize = usernameSize
		}
	}

	format := fmt.Sprintf("%%-%ds: %%d\n", usernameMaxSize)

	for user, score := range p.status.Score {
		fmt.Printf(format, user, score)
	}
}

func (p *Player) updateStatus() (err error) {
	p.status, err = p.Client.GetGameIdentifierStatus(p.status.Identifier)
	return
}

func (p *Player) startPlugins() {
	p.AIs.Start()
	p.Listeners.Start()
}

func brainNumber(a BasicAntStatus) int {
	if a.Brain == "controlled" {
		return 1
	}

	// we don't know what to expect here
	return 0
}

func visibilityNumber(c Cell) int {
	if c.Visibility {
		return 1
	}

	return 0
}

var contents = map[string]int{
	"grass": 0,
	"rock":  2,
	"water": 4,
	"sugar": 1,
	"mill":  3,
	"meat":  5,
}

func contentNumber(c Cell) (v int) {
	// see the format spec
	v, ok := contents[c.Content]
	if !ok {
		v = 0
	}

	return
}

func (p *Player) sendTurnStatusToPlugins() {
	var buf bytes.Buffer

	playing := 1
	if p.status.Status == "over" {
		p.done = true
		playing = 0
	}

	otherAnts := make(map[Position]BasicAntStatus)

	buf.WriteString(fmt.Sprintf("%d %d %d %d\n",
		p.turn.Number,                    // T
		p.status.Game.Spec.AntsPerPlayer, // A
		len(p.status.Players),            // P
		playing,                          // S
	))

	for _, ant := range p.turn.AntsStatuses {
		// save other visible ants
		for _, other := range ant.OtherVisibleAnts() {
			otherAnts[other.Pos] = other
		}

		// update the current map
		ant.Vision.SetVisibility(true)
		p.partialMap.ResetVisibility()
		p.partialMap.Combine(*ant.Vision)

		buf.WriteString(fmt.Sprintf("%d %d %d %d %d %d %d %d\n",
			ant.ID,                          // ID
			ant.Pos.X,                       // X
			ant.Pos.Y,                       // Y
			ant.Dir.X,                       // DX
			ant.Dir.Y,                       // DY
			ant.Energy,                      // E
			ant.Acid,                        // A
			brainNumber(ant.BasicAntStatus), // B
		))
	}

	// N
	buf.WriteString(fmt.Sprintf("%d\n", len(otherAnts)))

	for _, ant := range otherAnts {
		buf.WriteString(fmt.Sprintf("%d %d %d %d %d\n",
			ant.Pos.X,        // X
			ant.Pos.Y,        // Y
			ant.Dir.X,        // DX
			ant.Dir.Y,        // DY
			brainNumber(ant), // B
		))
	}

	buf.WriteString(fmt.Sprintf("%d %d %d\n",
		p.partialMap.Width(),    // W
		p.partialMap.Height(),   // H
		len(p.partialMap.Cells), // N
	))

	for _, cell := range p.partialMap.Cells {
		buf.WriteString(fmt.Sprintf("%d %d %d %d\n",
			cell.Pos.X,              // X
			cell.Pos.Y,              // Y
			contentNumber(*cell),    // C
			visibilityNumber(*cell), // S
		))
	}

	msg := buf.String()

	if p.debug {
		fmt.Println(msg)
	}

	p.AIs.SendAll(msg)
	p.Listeners.SendAll(msg)
}

func (p *Player) playTurn() (err error) {
	cmd := p.AIs.GetCommandResponse()

	p.turn, err = p.Client.PlayIdentifier(p.status.Identifier, cmd)
	if err != nil {
		return
	}

	err = p.updateStatus()
	return
}
