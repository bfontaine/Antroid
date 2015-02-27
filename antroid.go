package main

import (
	"fmt"
	"github.com/bfontaine/antroid/api"
	"os"
)

func main() {
	// just a demo
	h := api.NewHTTClient()
	res, err := h.CallApi()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", res)
}
