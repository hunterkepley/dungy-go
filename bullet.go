package main

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten"
)

// BulletExplosion is the animation that plays when a bullet hits something
type BulletExplosion struct {
	position Vec2f
	size     Vec2i

	animation         Animation
	animationSpeed    float64
	animationStarted  bool
	animationFinished bool

	sprite *Sprite

	image *ebiten.Image
}

func createBulletExplosion(position Vec2f, image *ebiten.Image) BulletExplosion {
	spritesheet := createSpritesheet(
		newVec2i(0, 60),
		newVec2i(112, 76),
		7,
		iitemsSpritesheet,
	)
	size := spritesheet.sprites[0].size
	sprite := createSprite(newVec2i(0, 60), newVec2i(112, 76), size, image)
	position.x -= float64(size.x) / 2
	position.y -= float64(size.y) / 2
	animation := createAnimation(
		spritesheet,
		image,
	)

	return BulletExplosion{
		position: position,
		size:     size,

		animation:         animation,
		animationSpeed:    6.,
		animationStarted:  false,
		animationFinished: false,

		sprite: &sprite,

		image: image,
	}
}

func (b *BulletExplosion) update(g *Game) {
	if !b.animationStarted {
		b.animation.startForwards()
		b.animationStarted = true
	}
	if b.animation.currentFrame == b.animation.spritesheet.numberOfSprites-1 {
		b.animationFinished = true
	}
	b.animation.update(b.animationSpeed)
}

func (b *BulletExplosion) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.position.x, b.position.y)
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?

	currentFrame := b.animation.spritesheet.sprites[b.animation.currentFrame]

	subImageRect := image.Rect(
		currentFrame.startPosition.x,
		currentFrame.startPosition.y,
		currentFrame.endPosition.x,
		currentFrame.endPosition.y,
	)

	screen.DrawImage(b.image.SubImage(subImageRect).(*ebiten.Image), op)
}

func updateBulletExplosions(game *Game) {
	for i := 0; i < len(game.bulletExplosions); i++ {
		if i-1 > len(game.bulletExplosions) {
			break
		}

		game.bulletExplosions[i].update(game)

		if game.bulletExplosions[i].animationFinished && game.bulletExplosions[i].animation.currentFrame == 0 {
			game.bulletExplosions = removeBulletExplosion(game.bulletExplosions, i)
		}
	}
}

func renderBulletExplosions(game *Game, screen *ebiten.Image) {
	for i := 0; i < len(game.bulletExplosions); i++ {
		game.bulletExplosions[i].render(screen)
	}
}

func removeBulletExplosion(slice []BulletExplosion, e int) []BulletExplosion {
	return append(slice[:e], slice[e+1:]...)
}

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

	collisionRect image.Rectangle

	destroy bool

	glow       BulletGlow
	glowSprite *Sprite

	lightID int // The ID of the light that follows the bullet

	sprite Sprite

	image *ebiten.Image
}

func createBullet(position Vec2f, rotation float64, speed float64, lightID int) Bullet {

	velocity := newVec2f(speed*math.Cos(rotation), speed*math.Sin(rotation))

	glowSprite := createSprite(newVec2i(27, 46), newVec2i(34, 53), newVec2i(7, 7), iitemsSpritesheet)

	size := newVec2i(5, 5)

	return Bullet{
		position:   position,
		rotation:   rotation,
		sprite:     createSprite(newVec2i(22, 47), newVec2i(27, 52), size, iitemsSpritesheet),
		size:       size,
		image:      iitemsSpritesheet,
		velocity:   velocity,
		speed:      speed,
		glow:       BulletGlow{image: iitemsSpritesheet, sprite: &glowSprite},
		glowSprite: &glowSprite,
		lightID:    lightID,
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

func (b *Bullet) update(game *Game) {
	b.borderCollision()

	game.lightHandler.lights[game.lightHandler.getLightIndex(b.lightID)].
		update(
			newVec2f(b.position.x-float64(b.size.x/2), b.position.y-float64(b.size.y/2)),
		)

	b.position.x += b.velocity.x
	b.position.y += b.velocity.y

	// Move glow
	b.glow.update(b.position, b.rotation)

	b.collisionRect = image.Rect(
		int(b.position.x),
		int(b.position.y),
		int(b.position.x)+b.size.x,
		int(b.position.y)+b.size.y,
	)

}

// Checks if bullets collide with border, deletes if so
func (b *Bullet) borderCollision() {
	if b.position.x <= 17+float64(b.size.x) ||
		b.position.y <= 14 ||
		b.position.x+float64(b.size.x) >= screenWidth-17 ||
		b.position.y+float64(b.size.y) >= screenHeight-17 {

		b.destroy = true
	}
}
