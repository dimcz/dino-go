package main

import (
	"log"

	"github.com/dimcz/dino-go/game"
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	g, err := game.NewGame(1)
	if err != nil {
		log.Fatalf("failed to init: %s", err)
	}

	pixelgl.Run(g.Start)
}
