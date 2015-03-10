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

	cl.Logout()
}
