package main

import (
	"flag"
	"log"

	"github.com/dimcz/dino-go/game"
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	d := 1
	p := 10

	flag.IntVar(&d, "d", 10, "dinosaurs count")
	flag.IntVar(&p, "p", 10, "populations count")
	flag.Parse()

	g, err := game.NewGame(d, p)
	if err != nil {
		log.Fatalf("failed to init: %s", err)
	}

	pixelgl.Run(g.Start)
}
