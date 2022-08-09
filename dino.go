package main

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

const (
	RUN = iota
	JUMP
)

var colorScheme = map[string]color.RGBA{
	"default": colornames.Black,
}

type dino struct {
	sprites []*pixel.Sprite
	sprite  *pixel.Sprite

	name  string
	color string
	state int
	index int
}

func initDino(name, color string) (*dino, error) {
	d := dino{
		sprites: make([]*pixel.Sprite, 3),
		color:   color,
		name:    name,
		state:   RUN,
	}

	pic, err := loadPicture(fmt.Sprintf("assets/images/dino/%s_run1.png", color))
	if err != nil {
		return nil, err
	}
	d.sprites[0] = pixel.NewSprite(pic, pic.Bounds())

	pic, err = loadPicture(fmt.Sprintf("assets/images/dino/%s_run2.png", color))
	if err != nil {
		return nil, err
	}
	d.sprites[1] = pixel.NewSprite(pic, pic.Bounds())

	pic, err = loadPicture(fmt.Sprintf("assets/images/dino/%s_jump.png", color))
	if err != nil {
		return nil, err
	}
	d.sprites[2] = pixel.NewSprite(pic, pic.Bounds())

	d.sprite = d.sprites[0]

	return &d, nil
}

func (d *dino) update() {
	if d.state == RUN {
		d.sprite = d.sprites[d.index/5]
		d.index += 1
		if d.index >= 10 {
			d.index = 0
		}
	}
}

func (d *dino) draw(target *pixelgl.Window, atlas *text.Atlas) {
	d.update()

	shift := 100.0
	if d.state == JUMP {
		shift = 200.0
	}

	vec := pixel.V(
		d.sprite.Frame().W()/2+30,
		d.sprite.Frame().H()/2+shift,
	)
	d.sprite.Draw(target, pixel.IM.Moved(vec))

	txt := text.New(pixel.ZV, atlas)
	txt.Color = colorScheme[d.color]
	_, _ = txt.WriteString(d.name)

	vec = pixel.V(
		d.sprite.Frame().Max.X-d.sprite.Frame().W()/2-txt.Orig.X/2+30,
		d.sprite.Frame().Max.Y+shift+5)

	txt.Draw(target, pixel.IM.Moved(vec))
}
