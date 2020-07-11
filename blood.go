package main

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

// Blood is a blood drop
type Blood struct {
	position Vec2f
	rotation float64

	subImage image.Rectangle

	sprite Sprite
	image  *ebiten.Image
}

func (b *Blood) render(screen *ebiten.Image) {
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

func (b *Blood) update() {

}

// BloodEmitter sprays out blood onto the ground!
type BloodEmitter struct {
	position Vec2f

	image *ebiten.Image
}
