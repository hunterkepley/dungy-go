package main

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

// WormAnimations is the animations for worms
type WormAnimations struct {
	idleFront Animation
}

// WormAnimationSpeeds is the animation speeds for worms
type WormAnimationSpeeds struct {
	idle float64
}

// Worm is a worm type enemy
type Worm struct {
	position Vec2f

	spritesheet     Spritesheet
	animation       Animation
	animations      WormAnimations
	animationSpeeds WormAnimationSpeeds

	image *ebiten.Image
}

func createWorm(position Vec2f) *Worm {
	idleFrontSpritesheet := createSpritesheet(newVec2i(0, 0), newVec2i(68, 22), 4, ienemiesSpritesheet)

	return &Worm{
		position: position,

		spritesheet: idleFrontSpritesheet,
		animations: WormAnimations{
			idleFront: createAnimation(idleFrontSpritesheet, ienemiesSpritesheet),
		},
		animationSpeeds: WormAnimationSpeeds{
			idle: 0.9,
		},

		image: ienemiesSpritesheet,
	}
}

func (w *Worm) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(w.position.x, w.position.y)
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?
	subImageRect := image.Rect(
		w.spritesheet.sprites[w.animation.currentFrame].startPosition.x,
		w.spritesheet.sprites[w.animation.currentFrame].startPosition.y,
		w.spritesheet.sprites[w.animation.currentFrame].endPosition.x,
		w.spritesheet.sprites[w.animation.currentFrame].endPosition.y,
	)
	screen.DrawImage(w.image.SubImage(subImageRect).(*ebiten.Image), op) // Draw player
}

func (w *Worm) update(bullets []Bullet) {
	// Start the animation if it's not playing
	if w.animation.state != AnimationPlayingForwards {
		w.animation = w.animations.idleFront
		w.animation.startForwards()
	}
	w.animation.update(w.animationSpeeds.idle)

	for _, b := range bullets {

	}
}
