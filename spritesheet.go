package main

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

// Sprite stores information on where a sprite subimage is located
type Sprite struct {
	startPosition Vec2i
	endPosition   Vec2i
	size          Vec2i
}

func createSprite(startPosition Vec2i, endPosition Vec2i, size Vec2i, image *ebiten.Image) Sprite {
	return Sprite{startPosition, endPosition, size}
}

func (s *Sprite) getBounds() image.Rectangle {
	return image.Rect(s.startPosition.x, s.startPosition.y, s.endPosition.x, s.endPosition.y)
}

// Spritesheet is a collection of Sprite's
type Spritesheet struct {
	sprites         []Sprite
	numberOfSprites int
	startPosition   Vec2i
	size            Vec2i
}

func createSpritesheet(startPosition Vec2i, endPosition Vec2i, numberOfSprites int, image *ebiten.Image) Spritesheet {
	// Create sprite slice
	sprites := make([]Sprite, numberOfSprites)
	// Calculate size of entire sheet
	size := newVec2i(endPosition.x-startPosition.x, endPosition.y-startPosition.y)
	// Calculate size of single sprite
	spriteSize := newVec2i(size.x/numberOfSprites, size.y)
	for i := 0; i < numberOfSprites; i++ {
		// Calculate start and end positions of current sprite
		spriteStartPosition := newVec2i(
			i*spriteSize.x+startPosition.x,
			startPosition.y,
		)
		spriteEndPosition := newVec2i(
			i*spriteSize.x+startPosition.x+spriteSize.x,
			spriteSize.y+startPosition.y,
		)
		// Add sprite to slice
		sprites[i] = Sprite{spriteStartPosition, spriteEndPosition, spriteSize}
	}
	return Spritesheet{
		sprites,
		numberOfSprites,
		startPosition,
		size,
	}
}

func (s *Spritesheet) endPosition() Vec2i {
	return newVec2i(s.startPosition.x+s.size.x, s.startPosition.y+s.size.y)
}
