package main

import "github.com/hajimehoshi/ebiten"

var (
	// i_ prefix is for images
	i_playerSpritesheet *ebiten.Image
)

func loadPlayerImages() {
	i_playerSpritesheet, _ = loadImage("./Assets/Art/Player/spritesheet.png")
}
