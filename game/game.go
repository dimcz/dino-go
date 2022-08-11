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

type Game struct {
	dinosaurs Dinosaurs

	road    *road
	enemies *enemies

	fonts map[float64]*text.Atlas
	info  *text.Text

	ai bool
}

func NewGame() (*Game, error) {
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

	return &Game{
		fonts:   fonts,
		info:    info,
		road:    r,
		enemies: initEnemies(),
	}, nil
}

func (g *Game) Single() {
	d, err := newDino(g.fonts[smallSize], dinoNames[0], dinoColors[0], dinoPadding)
	if err != nil {
		log.Fatalf("error loading: %s", err)
	}

	g.dinosaurs = Dinosaurs{"Fluffy": d}
	pixelgl.Run(g.run)
}

func (g *Game) AI(count int) {
	g.ai = true

	g.dinosaurs = make(Dinosaurs, count)
	for i := 0; i < count; i++ {
		d, err := newDino(g.fonts[smallSize],
			dinoNames[i],
			dinoColors[i],
			float64n(dinoPadding-10, dinoPadding+100))
		if err != nil {
			log.Fatalf("error loading: %s", err)
		}

		g.dinosaurs[dinoNames[i]] = d
	}

	pixelgl.Run(g.run)
}

func (g *Game) reset() {
	g.road.reset()
	g.enemies.reset()
}

func (g *Game) run() {
	win := initScreen("Dino Game!")
	defer win.Destroy()

	frameTick := setFPS(gameFPS)
	defer frameTick.Stop()

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

		g.showInfo(win, fmt.Sprintf("Scores: %0.f\nSpeed: %0.f", math.Floor(score), gameSpeed))

		for k, d := range g.dinosaurs {
			if g.enemies.checkCollisions(d) {
				delete(g.dinosaurs, k)

				if !g.ai {
					g.endGame(win, k)

					break gameLoop
				}
			}
		}

		if len(g.dinosaurs) == 0 {
			win.SetClosed(true)

			break gameLoop
		}

		score += gameSpeed / 8
		if score > scoreSpeedUp {
			scoreSpeedUp += gameSpeed * 50
			gameSpeed += 1
		}

		for _, d := range g.dinosaurs {
			d.draw(win, gameSpeed)
		}

		g.road.draw(win, gameSpeed)
		g.enemies.draw(win, gameSpeed)

		win.Update()

		if frameTick != nil {
			<-frameTick.C
		}
	}
}

func (g *Game) showInfo(window *pixelgl.Window, txt string) {
	g.info.Clear()
	_, _ = g.info.WriteString(txt)
	vec := pixel.V(infoPadding, Height-g.info.LineHeight-infoPadding)
	g.info.Draw(window, pixel.IM.Moved(vec))
}

func (g *Game) endGame(window *pixelgl.Window, name string) {
	window.UpdateInput()

	txt := text.New(pixel.V(0, 0), g.fonts[bigSize])
	txt.Color = colornames.Red
	_, _ = fmt.Fprintf(txt, "%s DIED", name)

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
