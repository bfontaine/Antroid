package main

import (
	"fmt"
	"github.com/bfontaine/antroid/api"
	"os"
)

func main() {
	// just a demo
	h := api.NewHTTClient()
	res, err := h.Call(api.GET, api.CALL_API_INFO, api.NoParams)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", res)
}
