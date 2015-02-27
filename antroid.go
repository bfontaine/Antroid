package main

import (
	"fmt"
	"github.com/bfontaine/antroid/api"
	"os"
)

func exitErr(e error) {
	fmt.Printf("Error: %v\n", e)
	os.Exit(1)
}

func main() {
	// just a demo
	cl, _ := api.NewClient()

	if err := cl.LoginWithCredentials("ww", "a"); err != nil {
		exitErr(err)
	}

	s, err := cl.WhoAmI()

	if err != nil {
		exitErr(err)
	}

	fmt.Printf("username: %s\n", s)

	if err := cl.Logout(); err != nil {
		exitErr(err)
	}
}
