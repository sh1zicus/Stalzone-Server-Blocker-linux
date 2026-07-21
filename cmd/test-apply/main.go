package main

import (
	"log"

	"github.com/sh1zicus/stalzone-server-blocker/internal/nft"
	"github.com/sh1zicus/stalzone-server-blocker/internal/servers"
)

func main() {

	pools, err := servers.Load("data/Servers.json")
	if err != nil {
		log.Fatal(err)
	}

	pools[0].Tunnels[0].Selected = true

	if err := nft.Apply(pools); err != nil {
		log.Fatal(err)
	}

	log.Println("Rules applied")
}
