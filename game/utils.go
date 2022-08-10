package game

import (
	"image"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
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

func loadTTF(path string) (*truetype.Font, error) {
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

	return fnt, nil
}

func atlasTable(path string, sizes []float64) (map[float64]*text.Atlas, error) {
	fnt, err := loadTTF(path)
	if err != nil {
		return nil, err
	}

	ht := make(map[float64]*text.Atlas, len(sizes))
	for _, s := range sizes {
		face := truetype.NewFace(fnt, &truetype.Options{
			Size:              s,
			GlyphCacheEntries: 1,
		})
		ht[s] = text.NewAtlas(face, text.ASCII)
	}

	return ht, nil
}

func float64n(low, high float64) float64 {
	return low + rand.Float64()*(high-low)
}

func setFPS(fps int) *time.Ticker {
	if fps <= 0 {
		return nil
	}

	return time.NewTicker(time.Second / time.Duration(fps))
}
