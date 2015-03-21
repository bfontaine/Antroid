package api

// This file describes our local game server, implemented as a `Player` struct.

import (
	"bytes"
	"fmt"
	"os"
)

// A Player represents a local game server connected to the remote one and
// which sends and receive messages from it. It used an API client, an AI pool
// and a Listeners pool to communicate with the remote server, send turn infos
// to the AIs and get their commands back, and send turn infos to the plugins.
// One of these plugins can be the GUI, for example.
//
// See docs/ai_protocol.md for the protocol we use to communicate with AIs and
// Plugin Listeners.
//
// A player represents one player in one game. At each turn, it'll send a
// command to play followed by a call to get the game's status. The first time,
// it'll send a dummy play command where all ants are resting just to get their
// positions, since we don't have ants' positions when creating/joining a game.
type Player struct {
	// API client
	Client *Client
	// All AIs
	AIs *AIPool
	// All plugin Listeners
	Listeners *ListenersPool

	// User credentials
	username, password string

	// If this is true we'll print more info
	debug bool

	// The current game status
	status *GameStatus
	// The current turn
	turn *Turn
	// The map as we know it. This is updated at each turn
	partialMap *PartialMap

	// This will be true when the game will end
	done bool
}

// NewPlayer returns a pointer on a new Player
func NewPlayer(username, password string) (p *Player) {
	return &Player{
		Client:     NewClient(),
		AIs:        NewAIPool(),
		Listeners:  NewListenersPool(),
		username:   username,
		password:   password,
		status:     &GameStatus{},
		turn:       &EmptyTurn,
		partialMap: NewPartialMap(),
	}
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

	// we request the game's status to have all its parameters
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

	// "play" with this rest command
	p.turn, err = p.Client.PlayIdentifier(p.status.Identifier, commands)

	if err != nil {
		return
	}

	// start all plugins, including AIs
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

// Quit stops all AIs and Listeners and logout the player from the remote
// server
func (p *Player) Quit() error {
	p.AIs.Stop()
	p.Listeners.Stop()
	return p.Client.Logout()
}

// PrintScores prints the current scores
func (p *Player) PrintScores() {
	var usernameMaxSize int

	// get the longest username
	for username := range p.status.Score {
		usernameSize := len(username)

		if usernameSize > usernameMaxSize {
			usernameMaxSize = usernameSize
		}
	}

	// craft the output format to be large enough to contain the longest
	// username
	format := fmt.Sprintf("%%-%ds: %%d\n", usernameMaxSize)

	// print each score
	for user, score := range p.status.Score {
		fmt.Printf(format, user, score)
	}
}

// updateStatus calls the remote server for the game status and updates the
// player's one.
func (p *Player) updateStatus() (err error) {
	p.status, err = p.Client.GetGameIdentifierStatus(p.status.Identifier)
	return
}

// startPlugins starts all AIs and all Listeners
func (p *Player) startPlugins() {
	p.AIs.Start()
	p.Listeners.Start()
}

// internal helper to get a number for an ant's brain state
func brainNumber(a BasicAntStatus) int {
	if a.Brain == "controlled" {
		return 1
	}

	// we don't know what to expect here
	return 0
}

// internal helper to get a number for a cell's visibility
func visibilityNumber(c Cell) int {
	if c.Visibility {
		return 1
	}

	return 0
}

// cell contents, see `docs/ai_protocol.md`
var contents = map[string]int{
	"grass": 0,
	"rock":  2,
	"water": 4,
	"sugar": 1,
	"mill":  3,
	"meat":  5,
}

// internal helper to get a number for a cell's content
func contentNumber(c Cell) (v int) {
	v, ok := contents[c.Content]
	if !ok {
		v = 0
	}

	return
}

// constructs a message describing the current turn (see `docs/ai_protocol.md`)
// and send it to all plugins and AIs.
func (p *Player) sendTurnStatusToPlugins() {
	var buf bytes.Buffer

	playing := 1
	if p.status.Status == "over" {
		p.done = true
		playing = 0
	}

	// header
	buf.WriteString(fmt.Sprintf("%d %d %d %d\n",
		p.turn.Number,                    // T
		p.status.Game.Spec.AntsPerPlayer, // A
		len(p.status.Players),            // P
		playing,                          // S
	))

	var visibleAnts, enemyAnts []BasicAntStatus

	// all our ants
	for _, ant := range p.turn.AntsStatuses {
		// save visible ants
		for _, visible := range ant.VisibleAnts {
			visibleAnts = append(visibleAnts, visible)
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

	// construct enemyAnts as all visible ants minus our ones
	for _, visible := range visibleAnts {
		for _, ant := range p.turn.AntsStatuses {
			if visible == ant.BasicAntStatus {
				continue
			}
		}
		enemyAnts = append(enemyAnts, visible)
	}

	// N enemy ants
	buf.WriteString(fmt.Sprintf("%d\n", len(enemyAnts)))

	// enemy ants
	for _, ant := range enemyAnts {
		buf.WriteString(fmt.Sprintf("%d %d %d %d %d\n",
			ant.Pos.X,        // X
			ant.Pos.Y,        // Y
			ant.Dir.X,        // DX
			ant.Dir.Y,        // DY
			brainNumber(ant), // B
		))
	}

	// map header
	buf.WriteString(fmt.Sprintf("%d %d %d\n",
		p.partialMap.Width(),    // W
		p.partialMap.Height(),   // H
		len(p.partialMap.Cells), // N
	))

	// map cells
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
		// print the message if we're debugging
		fmt.Fprintf(os.Stderr, "%s\n", msg)
	}

	// send it
	p.AIs.SendAll(msg)
	p.Listeners.SendAll(msg)
}

// playTurn gets the command to use from all AIs, send it to the server and
// updates the local game status.
func (p *Player) playTurn() (err error) {
	cmd := p.AIs.GetCommandResponse()

	p.turn, err = p.Client.PlayIdentifier(p.status.Identifier, cmd)
	if err != nil {
		return
	}

	err = p.updateStatus()
	return
}
