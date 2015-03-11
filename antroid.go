package main

import (
	"flag"
	"fmt"
	"github.com/bfontaine/antroid/api"
	"os"
	"strings"
)

// just for tests
const (
	username = "ww"
	password = "a"
)

func exitErr(e error) {
	fmt.Printf("Error: %v\n", e)
	os.Exit(1)
}

func main() {
	apiMethods := flag.Bool("api-methods", false, "show all API methods")
	whoAmI := flag.Bool("whoami", false, "verify that the user is logged")
	gamesList := flag.Bool("games", false, "list all the visible games")
	gameStatusId := flag.String("game-status", "", "get a game's status")
	destroyId := flag.String("destroy", "", "destroy a game")

	createGame := flag.Bool("create", false, "create a game")

	gs := api.GameSpec{Public: true}

	flag.StringVar(&gs.Description, "description", "", "game description")
	flag.IntVar(&gs.Pace, "pace", 1, "pace")
	flag.IntVar(&gs.Turns, "turns", 1, "turns")
	flag.IntVar(&gs.AntsPerPlayer, "ants", 1, "ants per player")
	flag.IntVar(&gs.MaxPlayers, "max", 1, "max players")
	flag.IntVar(&gs.MinPlayers, "min", 1, "min players")
	flag.IntVar(&gs.InitialEnergy, "energy", 100, "initial energy")
	flag.IntVar(&gs.InitialAcid, "acid", 100, "initial acid")

	flag.Parse()

	cl, _ := api.NewClient()

	if err := cl.LoginWithCredentials("ww", "a"); err != nil {
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

	if *gameStatusId != "" {
		gId := api.GameID(*gameStatusId)

		if gs, err := cl.GetGameIdentifierStatus(gId); err != nil {
			exitErr(err)
		} else {
			fmt.Printf("%s\n", gs)
		}
	}

	if *destroyId != "" {
		gId := api.GameID(*destroyId)

		if err := cl.DestroyGameIdentifier(gId); err != nil {
			exitErr(err)
		} else {
			fmt.Printf("Game %s successfully destroyed\n", gId)
		}
	}

	if *createGame {
		if g, err := cl.CreateGame(&gs); err != nil {
			exitErr(err)
		} else {
			fmt.Printf("Game %s successfully created\n", g.Identifier)
		}

	}

	cl.Logout()
}
