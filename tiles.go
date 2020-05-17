package main

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

const (
	wallOffset = 6 // How big the walls are, start tiles at this
)

// TileType is a type for the tyletype enum
type TileType int

const (
	// SmallTile ... TILETYPE ENUM [1]
	SmallTile TileType = iota + 1
	// BigTile ... TILETYPE ENUM [2]
	BigTile
	// WallTile ... TILETYPE ENUM [3]
	WallTile
)

func (t TileType) String() string {
	return [...]string{"Unknown", "SmallTile", "BigTile", "WallTile"}[t]
}

// ^ TILETYPE ENUM ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Tile be the tiles in the game
type Tile struct {
	position Vec2f
	size     Vec2i
	tileType TileType

	sprite Sprite
	image  *ebiten.Image // Spritesheet
}

var (
	smallTileSize = newVec2i(16, 17)
	bigTileSize   = newVec2i(31, 32)
	wallTileSize  = newVec2i(16, 29)
)

func createTile(position Vec2f, tileType TileType, image *ebiten.Image) Tile {
	var sprite Sprite
	size := smallTileSize
	switch tileType {
	case (SmallTile):
		sprite = createSprite(newVec2i(0, 0), smallTileSize, smallTileSize, image)
		break
	case (BigTile):
		sprite = createSprite(newVec2i(0, 18), newVec2i(31, 50), bigTileSize, image)
		size = bigTileSize
		break
	case (WallTile):
		sprite = createSprite(newVec2i(0, 51), newVec2i(16, 80), wallTileSize, image)
		size = wallTileSize
		break
	}
	return Tile{
		position,
		size,
		tileType,

		sprite,
		image,
	}
}

func (t *Tile) update() {

}

func (t *Tile) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(t.position.x, t.position.y)
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?
	subImageRect := image.Rect(
		t.sprite.startPosition.x,
		t.sprite.startPosition.y,
		t.sprite.endPosition.x,
		t.sprite.endPosition.y,
	)
	screen.DrawImage(t.image.SubImage(subImageRect).(*ebiten.Image), op)
}

// Generate the wall tiles at the top of the screen
func generateWalls(image *ebiten.Image) []Tile {
	numberOfWalls := screenWidth / wallTileSize.x
	t := make([]Tile, numberOfWalls)
	for i := 0; i < numberOfWalls; i++ {
		t[i] = createTile(newVec2f(float64(i*wallTileSize.x), 0), WallTile, image)
	}
	return t
}
