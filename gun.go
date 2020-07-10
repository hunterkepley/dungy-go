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

	rotation    float64
	storedAngle float64
	flipped     bool

	bullets []Bullet

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

func (g *Gun) update(playerPosition Vec2f, cursorPosition Vec2i) {

	// Placement offset [circle]
	radius := 12.
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
	if !(angle > Pi/2.5 && angle < Pi/3+Pi/3) && !(angle < Pi/-3 && angle > Pi/-3-Pi/3) {
		g.position.x = playerPosition.x - radius*math.Cos(angle) // Starting position x
		g.position.y = playerPosition.y - radius*math.Sin(angle) // Starting position y

		g.storedAngle = angle
	} else {
		g.position.x = playerPosition.x - radius*math.Cos(g.storedAngle) // Starting position x
		g.position.y = playerPosition.y - radius*math.Sin(g.storedAngle) // Starting position y
	}

	// Make always face the mouse
	g.rotation = angle + Pi
}

// Creates the bullets n stuff
func (g *Gun) fire() {
	g.bullets = append(
		g.bullets,
		createBullet(g.position, g.rotation),
	)
}

func (g *Gun) renderBullets(screen *ebiten.Image) {
	for i := 0; i < len(g.bullets); i++ {
		g.bullets[i].render(screen)
	}
}

func (g *Gun) updateBullets() {
	fmt.Println(len(g.bullets))
	for i := 0; i < len(g.bullets); i++ {
		if g.bullets[i].destroy {
			removeBullet(g.bullets, i)
			//continue
		}
		g.bullets[i].update()
		if i-1 > len(g.bullets) {
			break
		}
	}
}

func removeBullet(slice []Bullet, e int) []Bullet {
	return append(slice[:e], slice[e+1:]...)
}
