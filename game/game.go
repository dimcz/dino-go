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

	ai bool
}

func (g *Game) Player() {
	d, err := initDino(dinoTypes[0], dinoPadding)
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
		t := dinoTypes[i]

		d, err := initDino(t, float64n(dinoPadding-10, dinoPadding+30))
		if err != nil {
			log.Fatalf("error loading: %s", err)
		}

		g.dinosaurs[t.Name] = d
	}

	pixelgl.Run(g.run)
}

func (g *Game) run() {
	win := initScreen("Dino Game!")

	frameTick := setFPS(60)

	at, err := atlasTable(fontPath, []float64{bigSize, normalSize, smallSize})
	if err != nil {
		log.Fatalf("error loading: %s", err)
	}

	infoText := text.New(pixel.V(0, 0), at[normalSize])
	infoText.Color = colornames.Black

	r := initRoad()
	e := initEnemies()

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

		printInfo(win, infoText, 10, Height-10,
			fmt.Sprintf("Scores: %0.f\nSpeed: %0.f", math.Floor(score), gameSpeed))

		for k, d := range g.dinosaurs {
			if e.checkCollisions(d) {
				delete(g.dinosaurs, k)

				if !g.ai {
					endGame(win, at[bigSize], k)

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
			if d.isActive {
				d.draw(win, at[smallSize], gameSpeed)
			}
		}

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
		log.Fatalf("cannot create new window: %s", err)
	}

	return win
}

func endGame(target *pixelgl.Window, atlas *text.Atlas, name string) {
	txt := text.New(pixel.V(0, 0), atlas)
	txt.Color = colornames.Red
	_, _ = fmt.Fprintf(txt, "%s DIED", name)

	vec := target.Bounds().Center().Sub(pixel.V(txt.Bounds().W()/2, txt.LineHeight/2))
	txt.Draw(target, pixel.IM.Moved(vec))

	target.Update()
	target.UpdateInputWait(0)
	target.SetClosed(true)
}
