package main

import (
	"flag"
	"log"

	"github.com/dimcz/dino-go/game"
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	dinosaurs := 1
	epoch := 10

	flag.IntVar(&dinosaurs, "d", 20, "dinosaurs count")
	flag.IntVar(&epoch, "e", 10, "populations count")
	flag.Parse()

	g, err := game.NewGame(dinosaurs, epoch)
	if err != nil {
		log.Fatalf("failed to init: %s", err)
	}

	pixelgl.Run(g.Start)
}
