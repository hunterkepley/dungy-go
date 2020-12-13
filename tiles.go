package main

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

/*this comment is here bc vscode go linter is trash and keeps making an extra / * every time I save*/

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
	rotation float64
	scale    Vec2f

	sprite    Sprite
	image     *ebiten.Image // Spritesheet
	imageRect image.Rectangle
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
		numberOfSmallTiles := 6
		randomStart := newVec2i(rand.Intn(numberOfSmallTiles)*17, 0)
		randomEnd := newVec2i(randomStart.x+smallTileSize.x, randomStart.y+smallTileSize.y)
		sprite = createSprite(randomStart, randomEnd, smallTileSize, image)
	case (BigTile):
		sprite = createSprite(newVec2i(0, 18), newVec2i(31, 50), bigTileSize, image)
		size = bigTileSize
	case (WallTile):
		sprite = createSprite(newVec2i(0, 51), newVec2i(16, 80), wallTileSize, image)
		size = wallTileSize
	}
	return Tile{
		position: position,
		size:     size,
		tileType: tileType,
		scale:    Vec2f{1, 1},

		sprite: sprite,
		image:  image,
	}
}

func (t *Tile) render(screen *ebiten.Image) {
	if t.tileType != Empty {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Rotate(t.rotation)
		op.GeoM.Scale(t.scale.x, t.scale.y)
		op.GeoM.Translate(t.position.x, t.position.y)
		op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?
		if t.imageRect.Empty() {
			// If empty, give rect based on sprite
			t.imageRect = image.Rect(
				t.sprite.startPosition.x,
				t.sprite.startPosition.y,
				t.sprite.endPosition.x,
				t.sprite.endPosition.y,
			)
		}
		screen.DrawImage(t.image.SubImage(t.imageRect).(*ebiten.Image), op)
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

func createWallsFromMap(m Map, image *ebiten.Image) [][]Tile {
	for i := 0; i < len(m.tiles); i++ {
		for j := 0; j < len(m.tiles[i]); j++ {
			if m.mapNodes[i][j] == 'x' {
				m.tiles[i][j] = createTile(m.tiles[i][j].position, WallTile, image) // Make a wall if not
			}
		}
	}
	return m.tiles
}

func getNumberOfTilesPossible() Vec2i {
	numberOfTiles = newVec2i(screenHeight/smallTileSize.y, screenWidth/smallTileSize.x)
	return numberOfTiles
}

// Generate the tiles for the game
func generateTiles(image *ebiten.Image) [][]Tile {
	numberOfTiles := getNumberOfTilesPossible()
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

func renderTiles(g *Game, screen *ebiten.Image) {
	for _, w := range g.walls {
		w.render(screen)
	}
	for i := 0; i < len(g.currentMap.tiles); i++ {
		for j := 0; j < len(g.currentMap.tiles[i]); j++ {
			g.currentMap.tiles[i][j].render(screen)
		}
	}
}

func getRandomTile(g *Game, lessen Vec2i) (Tile, Vec2i) {
	index := newVec2i(rand.Intn(len(g.currentMap.tiles)-lessen.x), rand.Intn(len(g.currentMap.tiles[0])-lessen.y))

	return g.currentMap.tiles[index.x][index.y], index
}
