package main

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten"
)

// Bullet is for the player's gun to fire
type Bullet struct {
	position Vec2f
	size     Vec2i

	velocity    Vec2f
	rotation    float64
	storedAngle float64
	speed       float64

	sprite Sprite

	image *ebiten.Image
}

func createBullet(position Vec2f, rotation float64) Bullet {

	speed := 4.

	velocity := newVec2f(speed*math.Cos(rotation), speed*math.Sin(rotation))

	return Bullet{
		position: position,
		rotation: rotation,
		sprite:   createSprite(newVec2i(0, 14), newVec2i(6, 17), newVec2i(6, 3), iitemsSpritesheet),
		image:    iitemsSpritesheet,
		velocity: velocity,
		speed:    speed,
	}
}

func (b *Bullet) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// Translate center of image to 0, 0 before rotating
	op.GeoM.Translate(0-float64(b.sprite.size.x)/2, 0-float64(b.sprite.size.y)/2)
	op.GeoM.Rotate(b.rotation)
	op.GeoM.Translate(b.position.x, b.position.y)
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?

	subImageRect := image.Rect(
		b.sprite.startPosition.x,
		b.sprite.startPosition.y,
		b.sprite.endPosition.x,
		b.sprite.endPosition.y,
	)

	screen.DrawImage(b.image.SubImage(subImageRect).(*ebiten.Image), op)
}

func (b *Bullet) update() {
	b.position.x += b.velocity.x
	b.position.y += b.velocity.y
}

// Checks if bullets collide with border, deletes if so
func (b *Bullet) borderCollision() {

}
