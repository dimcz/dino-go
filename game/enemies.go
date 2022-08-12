package game

import (
	"github.com/faiface/pixel/pixelgl"
)

type enemies struct {
	cacti []*cactus
}

func distance(start float64) float64 {
	for {
		x := start + Width/float64n(0.8, 3)
		if x >= Width {
			return x
		}
	}
}

func initEnemies() *enemies {
	cacti := make([]*cactus, maxCountOfCacti)
	x := Width
	for i := 0; i < maxCountOfCacti; i++ {
		x = distance(x)
		cacti[i] = initCactus(x)
	}

	return &enemies{cacti}
}

func (e *enemies) draw(target *pixelgl.Window, step float64) (removed bool) {
	if e.cacti[0].x+e.cacti[0].sprite.Frame().W() <= 0 {
		e.cacti = e.cacti[1:]
		e.cacti = append(e.cacti, initCactus(distance(e.cacti[1].x)))
		removed = true
	}

	for _, v := range e.cacti {
		v.draw(target)
		v.x -= step
	}

	return removed
}

func (e *enemies) checkCollisions(d *dino) bool {
	for _, c := range e.cacti {
		if c.actual().Intersects(d.actual()) {
			return true
		}
	}

	return false
}
