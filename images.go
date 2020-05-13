package main

import "github.com/hajimehoshi/ebiten"

var (
	gameSpritesheet *ebiten.Image
)

func loadPlayerImages() {
	gameSpritesheet, _ = loadImage("./Assets/Art/gameSpritesheet.png")
}
