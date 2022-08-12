package game

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/yaricom/goNEAT/v3/neat/genetics"
	"golang.org/x/image/colornames"
)

const (
	RUN = iota
	JUMP
)

var colorScheme = map[string]color.RGBA{
	"default": colornames.Black,
	"aqua":    colornames.Aqua,
	"black":   colornames.Black,
	"bloody":  colornames.Red,
	"cobalt":  colornames.Aquamarine,
	"gold":    colornames.Gold,
	"insta":   colornames.Blueviolet,
	"lime":    colornames.Lime,
	"magenta": colornames.Magenta,
	"magma":   colornames.Darkred,
	"navy":    colornames.Navy,
	"neon":    colornames.Violet,
	"orange":  colornames.Orange,
	"pinky":   colornames.Pink,
	"purple":  colornames.Purple,
	"rgb":     colornames.Tomato,
	"silver":  colornames.Silver,
	"subaru":  colornames.Teal,
	"sunny":   colornames.Springgreen,
	"toxic":   colornames.Sandybrown,
}

type dino struct {
	sprites []*pixel.Sprite
	sprite  *pixel.Sprite
	info    *text.Text
	org     *genetics.Organism

	x, y  float64
	power float64

	state int
	index int
}

func newDino(font *text.Atlas, name, color string, padding float64, org *genetics.Organism) (*dino, error) {
	d := dino{
		sprites: make([]*pixel.Sprite, 3),
		x:       padding,
		y:       runPosition,
		power:   jumpPower,
		state:   RUN,
		org:     org,
	}

	d.info = text.New(pixel.ZV, font)
	d.info.Color = colorScheme[color]
	_, _ = d.info.WriteString(name)

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

func (d *dino) draw(target *pixelgl.Window, gameSpeed float64) {
	d.update(gameSpeed)

	vec := pixel.V(
		d.sprite.Frame().W()/2+d.x,
		d.sprite.Frame().H()/2+d.y,
	)
	d.sprite.Draw(target, pixel.IM.Moved(vec))

	vec = pixel.V(
		d.sprite.Frame().Max.X-d.sprite.Frame().W()/2-d.info.Orig.X/2+d.x,
		d.sprite.Frame().Max.Y+d.y+textPadding)

	d.info.Draw(target, pixel.IM.Moved(vec))
}

func (d *dino) actual() pixel.Rect {
	return pixel.R(
		d.x,
		d.y,
		d.x+d.sprite.Frame().W(),
		d.y+d.sprite.Frame().W(),
	)
}
