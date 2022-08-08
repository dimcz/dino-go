package main

import (
	"image"
	"io/ioutil"
	"math/rand"
	"os"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return pixel.PictureDataFromImage(img), nil
}

func matrix(rect pixel.Rect, delta float64) pixel.Matrix {
	vec := pixel.V(delta+rect.W()/2, rect.H()/2+roadPosition)
	return pixel.IM.Moved(vec)
}

func float64n(low, high float64) float64 {
	return low + rand.Float64()*(high-low)
}

func loadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	fnt, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(fnt, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}
