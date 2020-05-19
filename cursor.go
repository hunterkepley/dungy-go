package main

import (
	"github.com/hajimehoshi/ebiten"
)

// Cursor is the cursor
type Cursor struct {
	position Vec2i
	size     Vec2i

	image *ebiten.Image
}

func createCursor(image *ebiten.Image) Cursor {
	sizeX, sizeY := image.Size()
	return Cursor{
		newVec2i(0, 0),
		newVec2i(sizeX, sizeY),

		image,
	}
}

func (c *Cursor) update() {
	x, y := ebiten.CursorPosition()
	c.position = newVec2i(x-c.size.x/2, y-c.size.y/2)
}

func (c *Cursor) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.position.x), float64(c.position.y))
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?

	screen.DrawImage(c.image, op)
}
