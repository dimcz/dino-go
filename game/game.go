package game

import (
	"fmt"
	"log"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/yaricom/goNEAT/v3/neat/genetics"
	"github.com/yaricom/goNEAT/v3/neat/network"
	"golang.org/x/image/colornames"
)

type Dinosaurs map[string]*dino

type Game struct {
	dinosaurs Dinosaurs

	road    *road
	enemies *enemies

	fonts map[float64]*text.Atlas
	info  *text.Text

	eID  int
	Stop bool
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

func NewNEATGame(org []*genetics.Organism, eID int) (*Game, error) {
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
	game.eID = eID

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
			g.Stop = true
			win.SetClosed(true)

			break
		case win.Pressed(pixelgl.KeySpace):
			for _, d := range g.dinosaurs {
				if d.state != JUMP {
					d.jump(gameSpeed)
				}
			}
		}

		s := fmt.Sprintf("Scores: %0.f\nSpeed: %0.f", math.Floor(score), gameSpeed)
		if g.eID > -1 {
			s = fmt.Sprintf("%s\nGeneration: %d", s, g.eID)
		}
		g.showInfo(win, s)

		for k, d := range g.dinosaurs {
			if g.eID > -1 {
				// control
				if getControl(d.org.Phenotype, d, g.enemies.cacti, gameSpeed) > 0.5 {
					d.jump(gameSpeed)
					d.org.Fitness -= 1
				}
			}

			// collisions
			if g.enemies.checkCollisions(d) {
				delete(g.dinosaurs, k)
				d.org.Fitness = -10
			}
		}

		if len(g.dinosaurs) == 0 {
			if g.eID < 0 {
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
		g.enemies.draw(win, gameSpeed)

		for _, d := range g.dinosaurs {
			d.draw(win, gameSpeed)

			if g.eID > -1 {
				d.org.Fitness += .5
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
		eID:     -1,
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

func getInputs(d *dino, c *cactus, speed float64) []float64 {
	dinoRect := d.actual()
	cactusRect := c.actual()

	dx := dinoRect.Max.X - cactusRect.Min.X + cactusRect.W()/2
	dy := dinoRect.Max.Y - cactusRect.Min.Y

	dist := math.Sqrt(dx*dx + dy*dy)

	return []float64{0.1, dinoRect.Max.Y, dist, cactusRect.W(), speed}
}

func getControl(net *network.Network, d *dino, cacti []*cactus, speed float64) float64 {
	c := cacti[0]
	for _, v := range cacti {
		if v.actual().Max.X > d.actual().Min.X {
			c = v
			break
		}
	}

	inputs := getInputs(d, c, speed)
	if err := net.LoadSensors(inputs); err != nil {
		log.Fatal(err)
	}

	if res, err := net.Activate(); !res || err != nil {
		log.Fatal(err)
	}

	return net.ReadOutputs()[0]
}
