package game

import (
	"fmt"
	"log"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var pics [6]pixel.Picture

func init() {
	var err error
	for i := 0; i < 6; i++ {
		pics[i], err = loadPicture(fmt.Sprintf("assets/images/cactus/%d.png", i+1))
		if err != nil {
			log.Fatalf("error loading: %s", err)
		}
	}
}

type cactus struct {
	sprite *pixel.Sprite
	x      float64
}

func initCactus(x float64) *cactus {
	pic := pics[0] // rand.Intn(6)]
	sprite := pixel.NewSprite(pic, pic.Bounds())

	return &cactus{
		sprite: sprite,
		x:      x,
	}
}

func (c *cactus) draw(target *pixelgl.Window) {
	vec := pixel.V(
		c.x+c.sprite.Frame().W()/2,
		cactusPosition+c.sprite.Frame().H()/2,
	)
	c.sprite.Draw(target, pixel.IM.Moved(vec))
}

func (c *cactus) actual() pixel.Rect {
	return pixel.R(
		c.x,
		cactusPosition,
		c.x+c.sprite.Frame().W(),
		cactusPosition+c.sprite.Frame().H(),
	)
}
