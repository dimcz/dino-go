package main

import (
	"log"

	"github.com/dimcz/dino-go/game"
)

func main() {
	g, err := game.NewGame()
	if err != nil {
		log.Fatalf("failed to init: %s", err)
	}

	g.Run()
}
