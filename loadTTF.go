package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

func loadTTF(path string, dpi float64, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}

// InitFonts initializes fonts for the game
func (g *Game) InitFonts() {

	const dpi = 72
	var err error
	mdataFont, err = loadTTF("./Assets/Font/LoRe.ttf", dpi, 8)
	if err != nil {
		log.Fatal(err)
	}
	mversionFont, err = loadTTF("./Assets/Font/LoRe.ttf", dpi, 8)
	if err != nil {
		log.Fatal(err)
	}
}
