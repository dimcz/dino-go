package main

import (
	"log"

	"github.com/dimcz/dino-go/game"
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	dinosaurs := 1
	g, err := game.NewGame(dinosaurs)
	if err != nil {
		log.Fatalf("failed to init: %s", err)
	}

	pixelgl.Run(g.Start)
}
