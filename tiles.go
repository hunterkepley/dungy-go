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
	// SmallTile ... TILETYPE ENUM [0]
	SmallTile TileType = iota
	// BigTile ... TILETYPE ENUM [1]
	BigTile
	// WallTile ... TILETYPE ENUM [2]
	WallTile
)

func (t TileType) String() string {
	return [...]string{"SmallTile", "BigTile", "WallTile"}[t]
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

func createTile(position Vec2f, tileType TileType) Tile {
	size := newVec2i(0, 0)
	var sprite Sprite
	image := itileSpritesheet
	switch tileType {
	case (SmallTile):
		size = newVec2i(16, 17)
		sprite = createSprite(newVec2i(0, 0), size, size, image)
		break
	case (BigTile):
		size = newVec2i(31, 32)
		sprite = createSprite(newVec2i(0, 18), newVec2i(31, 50), size, image)
		break
	case (WallTile):
		size = newVec2i(16, 29)
		sprite = createSprite(newVec2i(0, 51), newVec2i(16, 80), size, image)
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
