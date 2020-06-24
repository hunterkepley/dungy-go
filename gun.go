package main

import (
	"fmt"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten"
)

// Gun is the player's gun
type Gun struct {
	position Vec2f

	rotation float64
	flipped  bool

	sprite Sprite
	image  *ebiten.Image
}

func (g *Gun) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// Translate center of image to 0, 0 before rotating
	op.GeoM.Translate(0-float64(g.sprite.size.x)/2, 0-float64(g.sprite.size.y)/2)
	if g.flipped {
		op.GeoM.Scale(1, -1)
	} else {
		op.GeoM.Scale(1, 1)
	}
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
	radius := 15.
	angle := math.Atan2(playerPosition.y-float64(cursorPosition.y), playerPosition.x-float64(cursorPosition.x))

	// Flip gun image
	if angle < Pi/2 && angle > Pi/-2 {
		if !g.flipped {
			g.flipped = true
			if angle > 0 {
				angle -= Pi / 6
			} else {
				angle += Pi / 6
			}
		}
	} else {
		if g.flipped {
			g.flipped = false
			if angle > 0 {
				angle += Pi / 6
			} else {
				angle -= Pi / 6
			}
		}
	}

	// This if statement only allows it to move on the sides of the body
	if !(angle > Pi/3 && angle < Pi/3+Pi/3) && !(angle < Pi/-3 && angle > Pi/-3-Pi/3) {
		g.position.x = playerPosition.x - radius*math.Cos(angle) // Starting position x
		g.position.y = playerPosition.y - radius*math.Sin(angle) // Starting position y

		fmt.Println(angle)
	}

	// TODO: Make the gun go to other side once it flips

	// Make always face the mouse
	g.rotation = angle + 135
}
