package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	roadImage    = "assets/images/road.png"
	roadPosition = 100
)

type road struct {
	pic         pixel.Picture
	left, right *pixel.Sprite
	pos         float64
}

func newRoad() *road {
	pic, err := loadPicture(roadImage)
	if err != nil {
		panic(err)
	}

	return &road{
		pic:   pic,
		left:  pixel.NewSprite(nil, pixel.ZR),
		right: pixel.NewSprite(nil, pixel.ZR),
	}
}

func (r *road) draw(target *pixelgl.Window, step float64) {
	if r.pos+Width > r.pic.Bounds().Max.X {
		x := r.pos + Width - r.pic.Bounds().Max.X
		if x < Width {
			rect := pixel.R(0, 0, x, r.pic.Bounds().Max.Y)
			r.right.Set(r.pic, rect)
			r.right.Draw(target, matrix(rect, Width-x))
		} else {
			r.pos = 0
		}
	}

	rect := pixel.R(r.pos, 0, r.pos+Width, r.pic.Bounds().Max.Y)
	r.left.Set(r.pic, rect)
	r.left.Draw(target, matrix(rect, 0))

	r.pos += step
}

func matrix(rect pixel.Rect, delta float64) pixel.Matrix {
	vec := pixel.V(delta+rect.W()/2, rect.H()/2+roadPosition)
	return pixel.IM.Moved(vec)
}
