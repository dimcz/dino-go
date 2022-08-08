package main

import (
	"github.com/faiface/pixel/pixelgl"
)

const (
	maxCountOfCacti = 3
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

func setEnemies() *enemies {
	cacti := make([]*cactus, maxCountOfCacti)
	x := Width
	for i := 0; i < maxCountOfCacti; i++ {
		x = distance(x)
		cacti[i] = newCactus(x)
	}

	return &enemies{cacti}
}

func (e *enemies) draw(target *pixelgl.Window, step float64) {
	if e.cacti[0].x+e.cacti[0].sprite.Frame().W() <= 0 {
		e.cacti = e.cacti[1:]
		e.cacti = append(e.cacti, newCactus(distance(e.cacti[1].x)))
	}

	for _, v := range e.cacti {
		v.draw(target)
		v.x -= step
	}
}
