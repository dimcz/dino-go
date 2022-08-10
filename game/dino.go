package game

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
	sprites  []*pixel.Sprite
	sprite   *pixel.Sprite
	dinoType Types

	x, y  float64
	power float64

	state int
	index int

	isActive bool
}

func initDino(dinoType Types, padding float64) (*dino, error) {
	d := dino{
		sprites:  make([]*pixel.Sprite, 3),
		dinoType: dinoType,
		x:        padding,
		y:        runPosition,
		power:    jumpPower,
		state:    RUN,
		isActive: true,
	}

	pic, err := loadPicture(fmt.Sprintf("assets/images/dino/%s_run1.png", dinoType.Color))
	if err != nil {
		return nil, err
	}
	d.sprites[0] = pixel.NewSprite(pic, pic.Bounds())

	pic, err = loadPicture(fmt.Sprintf("assets/images/dino/%s_run2.png", dinoType.Color))
	if err != nil {
		return nil, err
	}
	d.sprites[1] = pixel.NewSprite(pic, pic.Bounds())

	pic, err = loadPicture(fmt.Sprintf("assets/images/dino/%s_jump.png", dinoType.Color))
	if err != nil {
		return nil, err
	}
	d.sprites[2] = pixel.NewSprite(pic, pic.Bounds())

	d.sprite = d.sprites[0]

	return &d, nil
}

func (d *dino) reset() {
	d.y = runPosition
	d.power = jumpPower
	d.state = RUN
}

func (d *dino) update(gameSpeed float64) {
	if d.state == RUN {
		d.sprite = d.sprites[d.index/animationPause]
		d.index += 1
		if d.index >= animationPause*2 {
			d.index = 0
		}
	} else {
		d.jump(gameSpeed)
	}
}

func (d *dino) jump(gameSpeed float64) {
	if d.state == JUMP {
		d.y += d.power * gameSpeed / 4
		d.power -= gameSpeed / 16
		if d.y <= runPosition {
			d.y = runPosition
			d.state = RUN
			d.power = jumpPower
		}
	} else {
		d.state = JUMP
		d.sprite = d.sprites[2]
	}
}

func (d *dino) draw(target *pixelgl.Window, atlas *text.Atlas, gameSpeed float64) {
	d.update(gameSpeed)

	vec := pixel.V(
		d.sprite.Frame().W()/2+d.x,
		d.sprite.Frame().H()/2+d.y,
	)
	d.sprite.Draw(target, pixel.IM.Moved(vec))

	txt := text.New(pixel.ZV, atlas)
	txt.Color = colorScheme[d.dinoType.Color]
	_, _ = txt.WriteString(d.dinoType.Name)

	vec = pixel.V(
		d.sprite.Frame().Max.X-d.sprite.Frame().W()/2-txt.Orig.X/2+d.x,
		d.sprite.Frame().Max.Y+d.y+textPadding)

	txt.Draw(target, pixel.IM.Moved(vec))
}

func (d *dino) actual() pixel.Rect {
	return pixel.R(
		d.x,
		d.y,
		d.x+d.sprite.Frame().W(),
		d.y+d.sprite.Frame().W(),
	)
}
