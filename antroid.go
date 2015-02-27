package main

import (
	"fmt"
	"github.com/bfontaine/antroid/api"
	"os"
	"strings"
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

	info, err := cl.ApiInfo()

	if err != nil {
		exitErr(err)
	}

	var keys []string
	for k, _ := range info.Doc {
		keys = append(keys, k)
	}
	fmt.Printf("API methods: %s\n", strings.Join(keys, ", "))

	s, err := cl.WhoAmI()

	if err != nil {
		exitErr(err)
	}

	fmt.Printf("Username: %s\n", s)

	if err := cl.Logout(); err != nil {
		exitErr(err)
	}
}
