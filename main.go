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

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	pixelgl.Run(run)
}

func run() {
	win := initScreen("Dino Game!")

	frameTick := setFPS(60)

	atlasBig, err := newTTFAtlas("assets/fonts/intuitive.ttf", 100)
	if err != nil {
		panic(err)
	}

	atlas, err := newTTFAtlas("assets/fonts/intuitive.ttf", 26)
	if err != nil {
		panic(err)
	}

	atlasSmall, err := newTTFAtlas("assets/fonts/intuitive.ttf", 16)
	if err != nil {
		panic(err)
	}

	infoText := text.New(pixel.V(0, 0), atlas)
	infoText.Color = colornames.Black

	r := initRoad()
	e := initEnemies()

	d, err := initDino("Fluffy", "default")
	if err != nil {
		panic(err)
	}

	gameSpeed, score, scoreSpeedUp := 4.0, 0.0, 100.0

	isLoosing := false

	for !win.Closed() {
		win.Clear(colornames.White)

		switch {
		case win.JustPressed(pixelgl.KeyQ), win.Pressed(pixelgl.KeyEscape):
			win.SetClosed(true)
		case win.JustPressed(pixelgl.KeySpace):
			if d.state != JUMP {
				d.jump(gameSpeed)
			}
		case win.JustPressed(pixelgl.KeyR) && isLoosing:
			d.reset()
			e.reset()
			r.reset()
			gameSpeed, score, scoreSpeedUp = 4.0, 0.0, 100.0
			isLoosing = false
		}

		printInfo(win, infoText, 10, Height-10,
			fmt.Sprintf("Scores: %0.f\nSpeed: %0.f", math.Floor(score), gameSpeed))

		if isLoosing {
			txt := text.New(pixel.V(0, 0), atlasBig)
			txt.Color = colornames.Red
			_, _ = txt.WriteString("You DIED")

			vec := win.Bounds().Center().Sub(pixel.V(txt.Bounds().W()/2, txt.LineHeight/2))
			txt.Draw(win, pixel.IM.Moved(vec))
		} else {

			score += gameSpeed / 8
			if score > scoreSpeedUp {
				scoreSpeedUp += gameSpeed * 50
				gameSpeed += 1
			}

			d.draw(win, atlasSmall, gameSpeed)

			r.draw(win, gameSpeed)
			e.draw(win, gameSpeed)

			isLoosing = e.checkCollisions(d)
		}

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
