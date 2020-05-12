package main

import "github.com/hajimehoshi/ebiten"

var (
	playerSpritesheet *ebiten.Image
)

func loadPlayerImages() {
	playerSpritesheet, _ = loadImage("./Assets/Art/Player/spritesheet.png")
}
