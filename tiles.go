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
	numberOfWalls++ // One more wall could fit

	offset := newVec2f(10, 0) // Offset of the walls
	t := make([]Tile, numberOfWalls)
	for i := 0; i < numberOfWalls; i++ {
		// wallTileSize.x-1 to make them overlap on the x axis by 1 pixel
		t[i] = createTile(newVec2f(float64(i*(wallTileSize.x-1))+offset.x, offset.y), WallTile, image)
	}
	return t
}

// Generate the tiles for the game
func generateTiles(image *ebiten.Image) [][]Tile {
	numberOfTiles := newVec2i(screenWidth/smallTileSize.x, screenHeight/smallTileSize.y)
	numberOfTiles.x++ // One more tile could fit on x
	numberOfTiles.y--
	offset := newVec2f(10, float64(wallTileSize.y)) // Offset of the tiles
	t := [][]Tile{}
	for i := 0; i < numberOfTiles.x; i++ {
		t = append(t, []Tile{})
		for j := 0; j < numberOfTiles.y; j++ {
			// smallTileSize.x-1 to make them overlap on the x axis by 1 pixel
			t[i] = append(t[i], createTile(newVec2f(float64(i*(smallTileSize.x-1))+offset.x, float64(j*(smallTileSize.y-2))+offset.y), SmallTile, image))
		}
	}
	return t
}
