package main

import (
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

	fireSpeed    int
	firespeedMax int

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

func (g *Gun) update(playerPosition Vec2f, cursorCenter Vec2i) {

	// Placement offset [circle]
	radius := 12.
	angle := math.Atan2(playerPosition.y-float64(cursorCenter.y), playerPosition.x-float64(cursorCenter.x))

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

	if g.fireSpeed > 0 {
		g.fireSpeed--
	}

	// Make always face the mouse
	g.rotation = angle + Pi
}

// Creates the bullets n stuff
func (g *Gun) fire(playerPosition Vec2f, cursorCenter Vec2i) {
	if g.fireSpeed <= 0 {
		bulletSpeed := 0.2
		g.fireSpeed = g.firespeedMax
		g.bullets = append(
			g.bullets,
			createBullet(g.position, g.rotation, bulletSpeed),
		)
	}
}

func (g *Gun) renderBullets(screen *ebiten.Image) {
	for i := 0; i < len(g.bullets); i++ {
		g.bullets[i].render(screen)
	}
}

func (g *Gun) updateBullets() {
	for i := 0; i < len(g.bullets); i++ {
		// Break if some bullets were removed and for loop is too big
		if i-1 >= len(g.bullets) {
			break
		}
		g.bullets[i].update()

		// Destroy bullet if needed
		if g.bullets[i].destroy {
			g.bullets = removeBullet(g.bullets, i)
		}
	}
}

func removeBullet(slice []Bullet, e int) []Bullet {
	return append(slice[:e], slice[e+1:]...)
}
