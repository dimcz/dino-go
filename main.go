package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

const (
	Height = 720.0
	Width  = 1280.0
)

var (
	gameSpeed    = 4.0
	score        = 0.0
	scoreSpeedUp = 100.0
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	pixelgl.Run(run)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Dino Game!",
		Bounds: pixel.R(0, 0, Width, Height),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	face, err := loadTTF("assets/fonts/intuitive.ttf", 26)
	if err != nil {
		panic(err)
	}

	atlas := text.NewAtlas(face, text.ASCII)
	txt := text.New(pixel.V(0, 0), atlas)
	txt.Color = colornames.Black

	r := initRoad()
	e := initEnemies()

	for !win.Closed() {
		win.Clear(colornames.White)

		switch {
		case win.Pressed(pixelgl.KeyQ), win.Pressed(pixelgl.KeyEscape):
			win.SetClosed(true)
		}

		score += gameSpeed / 8
		if score > scoreSpeedUp {
			scoreSpeedUp += gameSpeed * 50
			gameSpeed += 1
		}

		printInfo(win, txt, 10, Height-10,
			fmt.Sprintf("Scores: %0.f\nSpeed: %0.f", math.Floor(score), gameSpeed))

		r.draw(win, gameSpeed)
		e.draw(win, gameSpeed)

		win.Update()
	}
}

func printInfo(win *pixelgl.Window, txt *text.Text, x, y float64, s string) {
	txt.Clear()
	_, _ = txt.WriteString(s)
	vec := pixel.V(x, y-txt.LineHeight)
	txt.Draw(win, pixel.IM.Moved(vec))
}
