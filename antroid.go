package main

// this executable is really a sandbox for now

import (
	"fmt"
	"github.com/bfontaine/antroid/api"
	"gopkg.in/alecthomas/kingpin.v1"
	"os"
	"strings"
)

func exitErr(e error) {
	fmt.Printf("Error: %v\n", e)
	os.Exit(1)
}

func gameServer(login, password string, ais []string, listeners []string,
	gs api.GameSpec, debug bool) {

	p := api.NewPlayer(login, password)

	p.SetDebug(debug)

	for _, ai := range ais {
		words := strings.Split(ai, " ")
		p.AIs.AddAI(words[0], words[1:]...)
	}

	for _, l := range listeners {
		words := strings.Split(l, " ")
		p.Listeners.AddListener(words[0], words[1:]...)
	}

	if err := p.Connect(); err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	if err := p.CreateAndJoinGame(&gs); err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	for {
		if done, err := p.PlayTurn(); err != nil {
			fmt.Printf("%s\n", err)
			return
		} else if done {
			fmt.Println("End of game.")
			fmt.Println("Scores:")
			p.PrintScores()
			return
		}
	}
}

var (
	app = kingpin.New("antroid", "A command-line Antroid API tool and game server")

	// global flags
	debug    = app.Flag("debug", "Enable debug mode.").Bool()
	login    = app.Flag("login", "Login.").Default("ww").String()
	password = app.Flag("password", "Password.").Default("a").String()

	// subcommands
	apiCmd     = app.Command("api", "Show all remote API methods.")
	whoCmd     = app.Command("whoami", "Show the logged user's name.")
	gamesCmd   = app.Command("games", "Show all visible games.")
	createCmd  = app.Command("create", "Create a new game.")
	statusCmd  = app.Command("status", "Get a game status.")
	destroyCmd = app.Command("destroy", "Destroy a game.")
	joinCmd    = app.Command("join", "Join a game.")
	playCmd    = app.Command("play", "Play a turn in a game.")
	serverCmd  = app.Command("server", "Start the local game server.")

	// play/server flags
	gameDesc = app.Flag("description", "Game description.").Default("a test").String()
	pace     = app.Flag("pace", "Game pace.").Default("1").Int()
	turns    = app.Flag("turns", "Number of turns.").Default("10").Int()
	ants     = app.Flag("ants", "Number of ants per player.").Default("1").Int()
	maxP     = app.Flag("max-players", "Max number of players.").Default("1").Int()
	minP     = app.Flag("min-players", "Min number of players.").Default("1").Int()
	energy   = app.Flag("energy", "Initial energy.").Default("100").Int()
	acid     = app.Flag("acid", "Initial acid.").Default("100").Int()
	players  = app.Flag("player", "Restrict games to this player "+
		"(can be used multiple times).").Strings()

	// subcommands args
	statusID  = statusCmd.Arg("id", "game ID").Required().String()
	destroyID = destroyCmd.Arg("id", "game ID").Required().String()
	joinID    = joinCmd.Arg("id", "game ID").Required().String()
	playID    = playCmd.Arg("id", "game ID").Required().String()
	playCmds  = playCmd.Arg("commands", "Commands to use for this turn.").Required().Strings()
	serverAIs = serverCmd.Arg("ais", "AIs to use for this game.").Required().Strings()

	// subcommands flags
	serverCreate = serverCmd.Flag("create", "Create a new game.").Bool()
	serverJoin   = serverCmd.Flag("join", "Join an existing game.").String()
	serverGui    = serverCmd.Flag("gui", "Use a GUI.").String()
)

func main() {
	app.Version("0.1.0")

	parsed := kingpin.MustParse(app.Parse(os.Args[1:]))

	gs := api.GameSpec{
		Description:   *gameDesc,
		Pace:          *pace,
		Turns:         *turns,
		AntsPerPlayer: *ants,
		MaxPlayers:    *maxP,
		MinPlayers:    *minP,
		InitialEnergy: *energy,
		InitialAcid:   *acid,
	}

	// construct the game spec
	if len(*players) == 0 {
		gs.Public = true
	} else {
		gs.Players = *players
	}

	if parsed == serverCmd.FullCommand() {
		if len(*serverAIs) == 0 {
			fmt.Fprintf(os.Stderr, "Expected at least one AI\n")
			os.Exit(1)
		}

		var plugins []string

		if *serverGui != "" {
			plugins = append(plugins, *serverGui)
		}

		gameServer(*login, *password, *serverAIs, plugins, gs, *debug)

		return
	}

	cl := api.NewClient()
	cl.SetDebug(*debug)

	if err := cl.LoginWithCredentials(*login, *password); err != nil {
		exitErr(err)
	}

	switch parsed {
	case apiCmd.FullCommand():
		if info, err := cl.APIInfo(); err != nil {
			exitErr(err)
		} else {
			var keys []string
			for k := range info.Doc {
				keys = append(keys, k)
			}
			fmt.Printf("API methods: %s\n", strings.Join(keys, ", "))
		}

	case whoCmd.FullCommand():
		if s, err := cl.WhoAmI(); err != nil {
			exitErr(err)
		} else {
			fmt.Printf("Username: %s\n", s)
		}

	case gamesCmd.FullCommand():
		if games, err := cl.ListGames(); err != nil {
			exitErr(err)
		} else {
			fmt.Println("Available games:")
			for _, g := range games {
				fmt.Printf("- %s\n", g)
			}
		}

	case statusCmd.FullCommand():
		gID := api.GameID(*statusID)

		if gs, err := cl.GetGameIdentifierStatus(gID); err != nil {
			exitErr(err)
		} else {
			fmt.Printf("%s\n", gs)
		}

	case destroyCmd.FullCommand():
		gID := api.GameID(*destroyID)

		if err := cl.DestroyGameIdentifier(gID); err != nil {
			exitErr(err)
		} else {
			fmt.Printf("Game %s successfully destroyed\n", gID)
		}

	case joinCmd.FullCommand():
		gID := api.GameID(*joinID)

		if err := cl.JoinGameIdentifier(gID); err != nil {
			exitErr(err)
		} else {
			fmt.Printf("Game %s successfully joined\n", gID)
		}

	case createCmd.FullCommand():
		if g, err := cl.CreateGame(&gs); err != nil {
			exitErr(err)
		} else {
			fmt.Printf("Game %s successfully created\n", g.Identifier)
		}

	case playCmd.FullCommand():
		cmds := strings.Join(*playCmds, ",")

		if t, err := cl.PlayIdentifier(api.GameID(*playID), api.Commands(cmds)); err != nil {
			exitErr(err)
		} else {
			fmt.Printf("%s\n", t)
		}
	default:
		app.Usage(os.Stderr)
	}

	cl.Logout()
}
