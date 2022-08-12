package game

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/yaricom/goNEAT/v3/neat/genetics"
	"github.com/yaricom/goNEAT/v3/neat/network"
	"golang.org/x/image/colornames"
)

type Dinosaurs map[string]*dino

type Game struct {
	count int
	age   int

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

func NewNEATGame(org []*genetics.Organism, age int) (*Game, error) {
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
	game.age = age

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
			os.Exit(0)
		case win.Pressed(pixelgl.KeySpace):
			for _, d := range g.dinosaurs {
				if d.state != JUMP {
					d.jump(gameSpeed)
				}
			}
		}

		g.showInfo(win, fmt.Sprintf("Scores: %0.f\nSpeed: %0.f\nGeneration: %d\n",
			math.Floor(score), gameSpeed, g.age+1))

		for k, d := range g.dinosaurs {
			if g.neat {
				// control
				if getControl(d.org.Phenotype, d, g.enemies.cacti, gameSpeed) > 0.2 {
					d.jump(gameSpeed)
					d.org.Fitness -= 0.5
				}
			}

			// collisions
			if g.enemies.checkCollisions(d) {

				if len(g.dinosaurs) > 1 {
					d.org.Fitness -= 1000
				}

				delete(g.dinosaurs, k)
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

		g.road.draw(win, gameSpeed)
		removed := g.enemies.draw(win, gameSpeed)

		for _, d := range g.dinosaurs {
			d.draw(win, gameSpeed)

			if g.neat && removed {
				d.org.Fitness += 1000
			}
		}

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

func getControl(net *network.Network, d *dino, cacti []*cactus, speed float64) float64 {
	c := cacti[0]
	for _, v := range cacti {
		if v.actual().Max.X > d.actual().Min.X {
			c = v
			break
		}
	}

	inputs := []float64{d.actual().Min.Y, getDistance(d.actual(), c.actual()), c.sprite.Frame().W(), speed}
	if err := net.LoadSensors(inputs); err != nil {
		log.Fatal(err)
	}

	if res, err := net.Activate(); !res || err != nil {
		log.Fatal(err)
	}

	return net.ReadOutputs()[0]
}

func getDistance(r1 pixel.Rect, r2 pixel.Rect) float64 {
	dx := r1.Max.X - r2.Min.X
	dy := r2.Max.Y - r2.Min.Y

	return math.Sqrt(dx*dx + dy*dy)
}
