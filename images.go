package main

import "github.com/hajimehoshi/ebiten"

var (
	// i prefix is for images
	iplayerSpritesheet  *ebiten.Image
	iinformationBox     *ebiten.Image
	itileSpritesheet    *ebiten.Image
	iUISpritesheet      *ebiten.Image
	icursor             *ebiten.Image
	iitemsSpritesheet   *ebiten.Image
	ienemiesSpritesheet *ebiten.Image
	ibloodSpritesheet   *ebiten.Image
)

func loadPlayerImages() {
	iplayerSpritesheet, _ = loadImage("./Assets/Art/Player/spritesheet.png")
}

func loadUIImages() {
	iUISpritesheet, _ = loadImage("./Assets/Art/UI/spritesheet.png")
}

func loadTileImages() {
	itileSpritesheet, _ = loadImage("./Assets/Art/Tiles/spritesheet.png")
}

func loadItemsImages() {
	iitemsSpritesheet, _ = loadImage("./Assets/Art/Items/spritesheet.png")
}

func loadEnemiesImages() {
	ienemiesSpritesheet, _ = loadImage("./Assets/Art/Enemies/spritesheet.png")
}

func loadParticlesImages() {
	ibloodSpritesheet, _ = loadImage("./Assets/Art/Particles/blood.png")
}
