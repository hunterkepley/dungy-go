package main

import "github.com/hajimehoshi/ebiten"

var (
	// i prefix is for images
	iplayerSpritesheet *ebiten.Image
	iinformationBox    *ebiten.Image
)

func loadPlayerImages() {
	iplayerSpritesheet, _ = loadImage("./Assets/Art/Player/spritesheet.png")
}

func loadUIImages() {
	iinformationBox, _ = loadImage("./Assets/Art/UI/informationBox.png")
}
