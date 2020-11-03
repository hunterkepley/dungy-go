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
	firing      bool // Is the gun firing? (Do we need to animate it?)

	fireSpeed    int
	firespeedMax int

	baseDamage int // Base damage the gun's bullets deliver

	animation      Animation // Fire animation
	animationSpeed float64   // Fire animation speed

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

	subImageRect := image.Rect(
		g.sprite.startPosition.x,
		g.sprite.startPosition.y,
		g.sprite.endPosition.x,
		g.sprite.endPosition.y,
	)

	if g.firing { // Gun is firing

		currentFrame := g.animation.spritesheet.sprites[g.animation.currentFrame]

		subImageRect = image.Rect(
			currentFrame.startPosition.x,
			currentFrame.startPosition.y,
			currentFrame.endPosition.x,
			currentFrame.endPosition.y,
		)

		flipAmount := 0.2

		if g.flipped {
			op.GeoM.Rotate(g.rotation + flipAmount)
		} else {
			op.GeoM.Rotate(g.rotation - flipAmount)
		}
	} else {
		op.GeoM.Rotate(g.rotation)
	}
	op.GeoM.Translate(g.position.x, g.position.y)

	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?

	screen.DrawImage(g.image.SubImage(subImageRect).(*ebiten.Image), op)
}

func (g *Gun) update(playerPosition Vec2f, cursorCenter Vec2i) {

	// Placement offset [circle]
	radius := 12.
	angle := math.Atan2(playerPosition.y-float64(cursorCenter.y), playerPosition.x-float64(cursorCenter.x))

	if g.firing {
		g.animation.update(g.animationSpeed)
		if g.animation.finishedFirstPlay {
			g.firing = false
			g.animation.finishedFirstPlay = false
		}
	}

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
func (g *Gun) fire(playerPosition Vec2f, game *Game) {
	if g.fireSpeed <= 0 {
		g.firing = true
		g.animation.startForwards()
		g.animation.currentFrame = 0

		bulletSpeed := 3.
		g.fireSpeed = g.firespeedMax
		lightID := game.lightHandler.addLight(game.lightHandler.lightImages.bulletLight, g.rotation)
		g.bullets = append(
			g.bullets,
			createBullet(g.position,
				g.rotation,
				bulletSpeed,
				lightID,
			),
		)
	}
}

func (g *Gun) renderBullets(screen *ebiten.Image) {
	for i := 0; i < len(g.bullets); i++ {
		g.bullets[i].render(screen)
	}
}

func (g *Gun) updateBullets(game *Game) {
	for i := 0; i < len(g.bullets); i++ {
		// Break if some bullets were removed and for loop is too big
		if i-1 >= len(g.bullets) {
			break
		}
		g.bullets[i].update(game)

		// Destroy bullet if needed
		if g.bullets[i].destroy {
			game.lightHandler.lights = removeLight(
				game.lightHandler.lights,
				g.bullets[i].lightID,
			)

			lightID := game.lightHandler.addLight(game.lightHandler.lightImages.bulletExplosionLight, g.rotation)

			game.bulletExplosions = append(
				game.bulletExplosions,
				createBulletExplosion(g.bullets[i].position, iitemsSpritesheet, lightID),
			)

			g.bullets = removeBullet(g.bullets, i)
		}
	}
}

func (g *Gun) calculateDamage() int {
	return g.baseDamage
}

func removeBullet(slice []Bullet, e int) []Bullet {
	removeLight(gameReference.lightHandler.lights, slice[e].lightID)
	return append(slice[:e], slice[e+1:]...)
}
