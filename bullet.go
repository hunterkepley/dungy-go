package main

import (
	"fmt"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten"
)

// BulletGlow is the glow around the laser bullets, just an image
type BulletGlow struct {
	position Vec2f
	size     Vec2i

	rotation float64

	sprite *Sprite

	image *ebiten.Image
}

func (b *BulletGlow) render(screen *ebiten.Image, glowSprite *Sprite) {
	op := &ebiten.DrawImageOptions{}
	// Translate center of image to 0, 0 before rotating
	op.GeoM.Translate(0-float64(b.size.x)/2, 0-float64(b.size.y)/2)
	op.GeoM.Rotate(b.rotation)
	op.GeoM.Translate(b.position.x, b.position.y)
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?

	fmt.Println("pos: ", b.position, "\nrot:", b.rotation, "\nsize:", b.size.x, ", ", b.size.y)

	subImageRect := image.Rect(
		b.sprite.startPosition.x,
		b.sprite.startPosition.y,
		b.sprite.endPosition.x,
		b.sprite.endPosition.y,
	)

	screen.DrawImage(b.image.SubImage(subImageRect).(*ebiten.Image), op)
}

func (b *BulletGlow) update(position Vec2f, rotation float64) {
	b.size = newVec2i(b.sprite.size.x, b.sprite.size.y)
	b.position = position
	b.rotation = rotation
}

// Bullet is for the player's gun to fire
type Bullet struct {
	position Vec2f
	size     Vec2i

	velocity    Vec2f
	rotation    float64
	storedAngle float64
	speed       float64

	destroy bool

	glow       BulletGlow
	glowSprite *Sprite

	sprite Sprite

	image *ebiten.Image
}

func createBullet(position Vec2f, rotation float64, speed float64) Bullet {

	velocity := newVec2f(speed*math.Cos(rotation), speed*math.Sin(rotation))

	glowSprite := createSprite(newVec2i(27, 46), newVec2i(34, 53), newVec2i(7, 7), iitemsSpritesheet)

	return Bullet{
		position:   position,
		rotation:   rotation,
		sprite:     createSprite(newVec2i(22, 47), newVec2i(27, 52), newVec2i(5, 5), iitemsSpritesheet),
		image:      iitemsSpritesheet,
		velocity:   velocity,
		speed:      speed,
		glow:       BulletGlow{image: iitemsSpritesheet, sprite: &glowSprite},
		glowSprite: &glowSprite,
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

	// Draw glow
	b.glow.render(screen, b.glowSprite)
}

func (b *Bullet) update() {
	if b.position.x <= 20+float64(b.size.x) ||
		b.position.y <= 14 ||
		b.position.x+float64(b.size.x) >= screenWidth-20-float64(b.size.x) ||
		b.position.y+float64(b.size.y) >= screenHeight-20-float64(b.size.y) {

		b.destroy = true
	}

	b.position.x += b.velocity.x
	b.position.y += b.velocity.y

	// Move glow
	b.glow.update(b.position, b.rotation)
}

// Checks if bullets collide with border, deletes if so
func (b *Bullet) borderCollision() {

}
