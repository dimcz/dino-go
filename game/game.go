package game

import (
	"fmt"
	"log"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

type Dinosaurs map[string]*dino

func (d Dinosaurs) notExists() bool {
	for _, d := range d {
		if d.isActive {
			return false
		}
	}

	return true
}

func (d Dinosaurs) draw(window *pixelgl.Window, step float64) {
	for _, d := range d {
		d.draw(window, step)
	}
}

type Game struct {
	count int

	dinosaurs Dinosaurs

	road    *road
	enemies *enemies

	fonts map[float64]*text.Atlas
	info  *text.Text
}

func NewGame(count int) (*Game, error) {
	fonts, err := atlasTable(fontPath, []float64{bigSize, normalSize, smallSize})
	if err != nil {
		return nil, err
	}

	info := text.New(pixel.V(0, 0), fonts[normalSize])
	info.Color = colornames.Black

	r, err := initRoad()
	if err != nil {
		return nil, err
	}

	dinosaurs := make(Dinosaurs, count)
	for i := 0; i < count; i++ {
		dinosaurs[dinoNames[i]], err = newDino(
			fonts[smallSize],
			dinoNames[i],
			dinoColors[i],
			float64n(dinoPadding-10, dinoPadding+100),
		)
		if err != nil {
			return nil, err
		}
	}

	return &Game{
		count:     count,
		fonts:     fonts,
		info:      info,
		road:      r,
		enemies:   initEnemies(),
		dinosaurs: dinosaurs,
	}, nil
}

func (g *Game) Start() {
	win := initScreen("Dino Game!")

	frameTick := setFPS(gameFPS)

	gameSpeed, score, scoreSpeedUp := 4.0, 0.0, 100.0

gameLoop:
	for !win.Closed() {
		win.Clear(colornames.White)

		switch {
		case win.Pressed(pixelgl.KeyQ), win.Pressed(pixelgl.KeyEscape):
			win.SetClosed(true)
		case win.Pressed(pixelgl.KeySpace):
			for _, d := range g.dinosaurs {
				if d.state != JUMP {
					d.jump(gameSpeed)
				}
			}
		}

		g.showInfo(win, fmt.Sprintf("Scores: %0.f\nSpeed: %0.f",
			math.Floor(score), gameSpeed))

		for _, d := range g.dinosaurs {
			if d.isActive && g.enemies.checkCollisions(d) {
				d.isActive = false
			}
		}

		if g.dinosaurs.notExists() {
			if g.count == 1 {
				g.endGame(win)
			}

			break gameLoop
		}

		score += gameSpeed / 8
		if score > scoreSpeedUp {
			scoreSpeedUp += gameSpeed * 50
			gameSpeed += 1
		}

		g.dinosaurs.draw(win, gameSpeed)
		g.road.draw(win, gameSpeed)
		g.enemies.draw(win, gameSpeed)

		win.Update()

		if frameTick != nil {
			<-frameTick.C
		}
	}

	g.reset()
}

func (g *Game) reset() {
	g.road.reset()
	g.enemies.reset()

	for _, d := range g.dinosaurs {
		d.reset()
	}
}

func (g *Game) showInfo(window *pixelgl.Window, txt string) {
	g.info.Clear()
	_, _ = g.info.WriteString(txt)
	vec := pixel.V(infoPadding, Height-g.info.LineHeight-infoPadding)
	g.info.Draw(window, pixel.IM.Moved(vec))
}

func (g *Game) endGame(window *pixelgl.Window) {
	window.UpdateInput()

	txt := text.New(pixel.V(0, 0), g.fonts[bigSize])
	txt.Color = colornames.Red
	_, _ = txt.WriteString("GAME OVER")

	vec := window.Bounds().Center().Sub(pixel.V(txt.Bounds().W()/2, txt.LineHeight/2))
	txt.Draw(window, pixel.IM.Moved(vec))

	window.Update()
	window.UpdateInputWait(0)

	window.SetClosed(true)
}

func initScreen(title string) *pixelgl.Window {
	cfg := pixelgl.WindowConfig{
		Title:   title,
		Bounds:  pixel.R(0, 0, Width, Height),
		Monitor: nil,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		log.Fatalf("cannot create new window: %s", err)
	}

	return win
}
