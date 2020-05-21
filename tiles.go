package main

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

const (
	wallOffset = 6 // How big the walls are, start tiles at this
)

var (
	numberOfTiles Vec2i
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
	// Empty ... TILETYPE ENUM [4]
	Empty // Used for holes or big tiles
)

func (t TileType) String() string {
	return [...]string{"Unknown", "SmallTile", "BigTile", "WallTile", "Empty"}[t]
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
	if t.tileType != Empty {
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
}

// Generate the wall tiles at the top of the screen
func generateWalls(image *ebiten.Image) []Tile {
	numberOfWalls := screenWidth / wallTileSize.x

	offset := newVec2f(17, 0) // Offset of the walls
	t := make([]Tile, numberOfWalls)
	for i := 0; i < numberOfWalls; i++ {
		// wallTileSize.x-1 to make them overlap on the x axis by 1 pixel
		t[i] = createTile(newVec2f(float64(i*(wallTileSize.x-1))+offset.x, offset.y), WallTile, image)
	}
	return t
}

// Generate the tiles for the game
func generateTiles(image *ebiten.Image) [][]Tile {
	numberOfTiles = newVec2i(screenHeight/smallTileSize.y, screenWidth/smallTileSize.x)
	numberOfTiles.x--
	offset := newVec2f(17, float64(wallTileSize.y)) // Offset of the tiles
	t := [][]Tile{}
	for i := 0; i < numberOfTiles.x; i++ {
		t = append(t, []Tile{})
		for j := 0; j < numberOfTiles.y; j++ {
			// smallTileSize.x-1 to make them overlap on the x axis by 1 pixel
			t[i] = append(t[i], createTile(newVec2f(float64(j*(smallTileSize.x-1))+offset.x, float64(i*(smallTileSize.y-2))+offset.y), SmallTile, image))
		}
	}
	return t
}

func generateBigTiles(tiles [][]Tile, image *ebiten.Image) {
	y := numberOfTiles.x - 2 // arrays start at 0, and for some reason it needs 1 more, so -2
	x := numberOfTiles.y - 2 // ^

	// middle
	makeBigTile(newVec2i(y/2, x/2), tiles, image)

	// top left tile
	makeBigTile(newVec2i(1, 1), tiles, image)
	// next to ^
	makeBigTile(newVec2i(1, 4), tiles, image)
	// below ^
	makeBigTile(newVec2i(4, 4), tiles, image)
	// next to ^
	makeBigTile(newVec2i(4, 7), tiles, image)

	// We have to flip X and Y because the way tiles were generated
	// bottom right tile
	makeBigTile(newVec2i(y-1, x-1), tiles, image)
	// next to ^
	makeBigTile(newVec2i(y-1, x-4), tiles, image)
	// above ^
	makeBigTile(newVec2i(y-4, x-4), tiles, image)
	// next to ^
	makeBigTile(newVec2i(y-4, x-7), tiles, image)

	// top right tile
	makeBigTile(newVec2i(1, x-1), tiles, image)
	// next to ^
	makeBigTile(newVec2i(1, x-4), tiles, image)
	// below ^
	makeBigTile(newVec2i(4, x-4), tiles, image)
	// next to ^
	makeBigTile(newVec2i(4, x-7), tiles, image)

	// bottom left tile
	makeBigTile(newVec2i(y-1, 1), tiles, image)
	// next to ^
	makeBigTile(newVec2i(y-1, 4), tiles, image)
	// above ^
	makeBigTile(newVec2i(y-4, 4), tiles, image)
	// next to ^
	makeBigTile(newVec2i(y-4, 7), tiles, image)
}

func makeBigTile(tilePosition Vec2i, tiles [][]Tile, image *ebiten.Image) {
	// Actual tile
	tiles[tilePosition.x][tilePosition.y] = createTile(
		tiles[tilePosition.x][tilePosition.y].position,
		BigTile,
		image,
	)
	// Make empties
	tiles[tilePosition.x+1][tilePosition.y].tileType = Empty // Don't need to waste time making a new tile
	tiles[tilePosition.x][tilePosition.y+1].tileType = Empty
	tiles[tilePosition.x+1][tilePosition.y+1].tileType = Empty
}
