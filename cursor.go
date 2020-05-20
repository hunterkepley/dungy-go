package main

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

// Cursor is the cursor
type Cursor struct {
	position Vec2i
	size     Vec2i

	currentCursor int
	cursors       []Sprite

	image *ebiten.Image
}

func createCursor(image *ebiten.Image) Cursor {
	return Cursor{
		newVec2i(0, 0),
		newVec2i(0, 0),

		0,
		[]Sprite{
			createSprite(newVec2i(0, 0), newVec2i(15, 15), newVec2i(15, 15), iUISpritesheet),
			createSprite(newVec2i(16, 0), newVec2i(30, 14), newVec2i(14, 14), iUISpritesheet),
			createSprite(newVec2i(31, 0), newVec2i(45, 14), newVec2i(14, 14), iUISpritesheet),
			createSprite(newVec2i(46, 0), newVec2i(50, 4), newVec2i(4, 4), iUISpritesheet),
		},

		image,
	}
}

func (c *Cursor) update() {
	x, y := ebiten.CursorPosition()
	if c.checkInScreen(x, y) {
		c.position = newVec2i(x-c.size.x/2, y-c.size.y/2)
	}
	c.size = newVec2i(c.cursors[c.currentCursor].size.x, c.cursors[c.currentCursor].size.y)
}

// Checks if the mouse is in the screen
func (c *Cursor) checkInScreen(x int, y int) bool {
	if x >= 0 && y >= 0 && x <= screenWidth && y <= screenHeight {
		return true
	}
	return false
}

func (c *Cursor) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.position.x), float64(c.position.y))
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?

	cursorRect := image.Rect(
		c.cursors[c.currentCursor].startPosition.x,
		c.cursors[c.currentCursor].startPosition.y,
		c.cursors[c.currentCursor].endPosition.x,
		c.cursors[c.currentCursor].endPosition.y,
	)
	screen.DrawImage(c.image.SubImage(cursorRect).(*ebiten.Image), op)
}
