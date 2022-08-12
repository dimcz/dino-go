package game

import (
	"fmt"
	"log"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/yaricom/goNEAT/v3/neat/genetics"
	"golang.org/x/image/colornames"
)

type Dinosaurs map[string]*dino

type Game struct {
	count int

	dinosaurs Dinosaurs

	road    *road
	enemies *enemies

	fonts map[float64]*text.Atlas
	info  *text.Text

	neat bool
}

func NewGame() (*Game, error) {
	game, err := newGame()
	if err != nil {
		return nil, err
	}

	d, err := newDino(
		game.fonts[smallSize],
		dinoNames[0],
		dinoColors[0],
		float64n(dinoPadding-10, dinoPadding+100),
		nil,
	)
	if err != nil {
		return nil, err
	}

	game.dinosaurs = Dinosaurs{dinoNames[0]: d}

	return game, nil
}

func NewNEATGame(org []*genetics.Organism) (*Game, error) {
	game, err := newGame()
	if err != nil {
		return nil, err
	}

	dinosaurs := make(Dinosaurs, len(org))
	for i, o := range org {
		dinosaurs[dinoNames[i]], err = newDino(
			game.fonts[smallSize],
			dinoNames[i],
			dinoColors[i],
			float64n(dinoPadding-10, dinoPadding+100),
			o,
		)
		if err != nil {
			return nil, err
		}
	}

	game.dinosaurs = dinosaurs
	game.neat = true

	return game, nil
}

func (g *Game) Run() {
	pixelgl.Run(g.start)
}

func (g *Game) start() {
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

		for k, d := range g.dinosaurs {
			if g.enemies.checkCollisions(d) {
				delete(g.dinosaurs, k)

				if g.neat {
					d.org.Fitness = -1
				}
			}
		}

		if len(g.dinosaurs) == 0 {
			if !g.neat {
				g.endGame(win)
			}

			break gameLoop
		}

		score += gameSpeed / 8
		if score > scoreSpeedUp {
			scoreSpeedUp += gameSpeed * 50
			gameSpeed += 1
		}

		for _, d := range g.dinosaurs {
			d.draw(win, gameSpeed)

			if g.neat {
				d.org.Fitness += 0.5
			}
		}

		g.road.draw(win, gameSpeed)
		g.enemies.draw(win, gameSpeed)

		win.Update()

		if frameTick != nil {
			<-frameTick.C
		}
	}
}

func newGame() (*Game, error) {
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
