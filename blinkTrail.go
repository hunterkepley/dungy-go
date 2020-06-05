package main

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

// BlinkTrail the trail when the player blinks
type BlinkTrail struct {
	speed    float64
	speedMax float64

	sections []BlinkTrailSection
}

// BlinkTrailSection is the actual animated section of the trail
type BlinkTrailSection struct {
	position Vec2f

	spritesheet Spritesheet
	animation   *Animation

	image *ebiten.Image
}

// BlinkTrail
func createBlinkTrail(speed float64) BlinkTrail {
	return BlinkTrail{
		speed:    speed,
		speedMax: speed,

		sections: []BlinkTrailSection{},
	}
}

func (b *BlinkTrail) render(screen *ebiten.Image) {
	for _, s := range b.sections {
		s.render(screen)
	}
}

func (b *BlinkTrail) update() {
	for _, s := range b.sections {
		s.update()
	}
}

func (b *BlinkTrail) spawnUpdate(position Vec2f) {
	// Add new section every few ticks
	if b.speed >= 0 {
		b.speed--
	} else {
		t := createBlinkTrailSection(position)
		b.sections = append(b.sections, t)
		b.speed = b.speedMax
		t.animation.startForwards()
	}
}

// BlinkTrailSection

func createBlinkTrailSection(position Vec2f) BlinkTrailSection {
	image := iplayerSpritesheet
	blinkSpritesheet := createSpritesheet(newVec2i(0, 207), newVec2i(39, 230), 3, image)
	animation := createAnimation(blinkSpritesheet, image)
	return BlinkTrailSection{
		position: position,

		spritesheet: blinkSpritesheet,
		animation:   &animation,

		image: image,
	}
}

func (b *BlinkTrailSection) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.position.x, b.position.y)
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?
	subImageRect := image.Rect(
		b.spritesheet.sprites[b.animation.currentFrame].startPosition.x,
		b.spritesheet.sprites[b.animation.currentFrame].startPosition.y,
		b.spritesheet.sprites[b.animation.currentFrame].endPosition.x,
		b.spritesheet.sprites[b.animation.currentFrame].endPosition.y,
	)
	screen.DrawImage(b.image.SubImage(subImageRect).(*ebiten.Image), op)
}

func (b *BlinkTrailSection) update() {
	b.animation.update(1.5)
}
