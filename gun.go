package main

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten"
)

// Gun is the player's gun
type Gun struct {
	position Vec2f
	offset   Vec2f

	rotation float64

	sprite Sprite
	image  *ebiten.Image
}

func (g *Gun) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// Translate center of image to 0, 0 before rotating
	op.GeoM.Translate(0-float64(g.sprite.size.x)/2, 0-float64(g.sprite.size.y)/2)
	op.GeoM.Rotate(g.rotation)
	op.GeoM.Translate(g.position.x, g.position.y)
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?

	subImageRect := image.Rect(
		g.sprite.startPosition.x,
		g.sprite.startPosition.y,
		g.sprite.endPosition.x,
		g.sprite.endPosition.y,
	)

	screen.DrawImage(g.image.SubImage(subImageRect).(*ebiten.Image), op)
}

func (g *Gun) update(playerPosition Vec2f) {

	cursorPosition := Vec2i{}
	cursorPosition.x, cursorPosition.y = ebiten.CursorPosition()

	// Placement offset [circle]
	g.offset = newVec2f(0, 0)
	radius := 20.
	angle := math.Atan2(playerPosition.y-float64(cursorPosition.y), playerPosition.x-float64(cursorPosition.x))

	g.position.x = playerPosition.x - radius*math.Cos(angle) // Starting position x
	g.position.y = playerPosition.y - radius*math.Sin(angle) // Starting position y

	// Make always face the mouse
	g.rotation = angle + 135
}
