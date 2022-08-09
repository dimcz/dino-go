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
	win := initScreen("Dino Game!")

	face, err := loadTTF("assets/fonts/intuitive.ttf", 26)
	if err != nil {
		panic(err)
	}

	frameTick := setFPS(60)

	atlas := text.NewAtlas(face, text.ASCII)
	txt := text.New(pixel.V(0, 0), atlas)
	txt.Color = colornames.Black

	faceSmall, err := loadTTF("assets/fonts/intuitive.ttf", 16)
	if err != nil {
		panic(err)
	}
	atlasSmall := text.NewAtlas(faceSmall, text.ASCII)

	r := initRoad()
	e := initEnemies()

	d, err := initDino("Fluffy", "default")
	if err != nil {
		panic(err)
	}

	for !win.Closed() {
		win.Clear(colornames.White)

		switch {
		case win.Pressed(pixelgl.KeyQ), win.Pressed(pixelgl.KeyEscape):
			win.SetClosed(true)
		case win.Pressed(pixelgl.KeySpace):
			if d.state != JUMP {
				d.jump()
			}
		}

		score += gameSpeed / 8
		if score > scoreSpeedUp {
			scoreSpeedUp += gameSpeed * 50
			gameSpeed += 1
		}

		printInfo(win, txt, 10, Height-10,
			fmt.Sprintf("Scores: %0.f\nSpeed: %0.f", math.Floor(score), gameSpeed))

		d.draw(win, atlasSmall)

		r.draw(win, gameSpeed)
		e.draw(win, gameSpeed)

		win.Update()

		if frameTick != nil {
			<-frameTick.C
		}
	}
}

func printInfo(win *pixelgl.Window, txt *text.Text, x, y float64, s string) {
	txt.Clear()
	_, _ = txt.WriteString(s)
	vec := pixel.V(x, y-txt.LineHeight)
	txt.Draw(win, pixel.IM.Moved(vec))
}

func initScreen(title string) *pixelgl.Window {
	cfg := pixelgl.WindowConfig{
		Title:   title,
		Bounds:  pixel.R(0, 0, Width, Height),
		Monitor: nil,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	return win
}

func setFPS(fps int) *time.Ticker {
	if fps <= 0 {
		return nil
	}

	return time.NewTicker(time.Second / time.Duration(fps))
}
