package main

// this executable is really a sandbox for now

import (
	"flag"
	"fmt"
	"github.com/bfontaine/antroid/api"
	"os"
	"strings"
)

func exitErr(e error) {
	fmt.Printf("Error: %v\n", e)
	os.Exit(1)
}

func gameServer(login, passwd, ai string, gs *api.GameSpec, debug bool) {
	p := api.NewPlayer(login, passwd)

	p.Client.SetDebug(debug)

	p.AIs.AddAI(ai)

	if err := p.Connect(); err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	if err := p.CreateAndJoinGame(gs); err != nil {
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

func main() {
	apiMethods := flag.Bool("methods", false, "show all API methods")
	whoAmI := flag.Bool("whoami", false, "verify that the user is logged")
	gamesList := flag.Bool("games", false, "list all the visible games")
	createGame := flag.Bool("create", false, "create a game")

	// flags with game ids
	gameStatusID := flag.String("status", "", "get a game's status")
	destroyID := flag.String("destroy", "", "destroy a game")
	joinID := flag.String("join", "", "join a game")
	playID := flag.String("play", "", "play a game")

	// -play options
	playCmds := flag.String("cmds", "", "play with this commands")
	prettyMap := flag.Bool("pretty", false, "print a pretty map when playing")

	gs := api.GameSpec{Public: true}

	// -create options
	flag.StringVar(&gs.Description, "description", "", "game description")
	flag.IntVar(&gs.Pace, "pace", 1, "pace")
	flag.IntVar(&gs.Turns, "turns", 10, "turns")
	flag.IntVar(&gs.AntsPerPlayer, "ants", 1, "ants per player")
	flag.IntVar(&gs.MaxPlayers, "max", 1, "max players")
	flag.IntVar(&gs.MinPlayers, "min", 1, "min players")
	flag.IntVar(&gs.InitialEnergy, "energy", 100, "initial energy")
	flag.IntVar(&gs.InitialAcid, "acid", 100, "initial acid")

	// general options
	login := flag.String("login", "ww", "login")
	passwd := flag.String("password", "a", "password")

	debug := flag.Bool("debug", false, "debug mode")

	server := flag.Bool("server", false, "start a server")

	ai := flag.String("ai", "", "the AI to use")

	flag.Parse()

	if *server {
		if *ai == "" {
			fmt.Println("AI expected")
			return
		}
		gameServer(*login, *passwd, *ai, &gs, *debug)
		return
	}

	cl := api.NewClient()

	cl.SetDebug(*debug)

	if err := cl.LoginWithCredentials(*login, *passwd); err != nil {
		exitErr(err)
	}

	if *apiMethods {
		if info, err := cl.APIInfo(); err != nil {
			exitErr(err)
		} else {
			var keys []string
			for k := range info.Doc {
				keys = append(keys, k)
			}
			fmt.Printf("API methods: %s\n", strings.Join(keys, ", "))
		}
	}

	if *whoAmI {
		if s, err := cl.WhoAmI(); err != nil {
			exitErr(err)
		} else {
			fmt.Printf("Username: %s\n", s)
		}
	}

	if *gamesList {
		if games, err := cl.ListGames(); err != nil {
			exitErr(err)
		} else {
			fmt.Println("Available games:")
			for _, g := range games {
				fmt.Printf("- %s\n", g)
			}
		}
	}

	if *gameStatusID != "" {
		gID := api.GameID(*gameStatusID)

		if gs, err := cl.GetGameIdentifierStatus(gID); err != nil {
			exitErr(err)
		} else {
			fmt.Printf("%s\n", gs)
		}
	}

	if *destroyID != "" {
		gID := api.GameID(*destroyID)

		if err := cl.DestroyGameIdentifier(gID); err != nil {
			exitErr(err)
		} else {
			fmt.Printf("Game %s successfully destroyed\n", gID)
		}
	}

	if *joinID != "" {
		gID := api.GameID(*joinID)

		if err := cl.JoinGameIdentifier(gID); err != nil {
			exitErr(err)
		} else {
			fmt.Printf("Game %s successfully joined\n", gID)
		}
	}

	if *createGame {
		if g, err := cl.CreateGame(&gs); err != nil {
			exitErr(err)
		} else {
			fmt.Printf("Game %s successfully created\n", g.Identifier)
		}
	}

	if *playID != "" {
		if t, err := cl.PlayIdentifier(api.GameID(*playID), api.Commands(*playCmds)); err != nil {
			exitErr(err)
		} else {
			if *prettyMap {
				fmt.Println(t.PrettyString())
			} else {
				fmt.Printf("%s\n", t)
			}
		}
	}

	cl.Logout()
}
