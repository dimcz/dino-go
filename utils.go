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

func float64n(low, high float64) float64 {
	return low + rand.Float64()*(high-low)
}
