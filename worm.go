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

	health    int
	maxHealth int

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

		health:    10,
		maxHealth: 10,

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
		size := newVec2i(
			w.spritesheet.sprites[w.animation.currentFrame].size.x,
			w.spritesheet.sprites[w.animation.currentFrame].size.y,
		)
		endPosition := newVec2i(
			int(w.position.x)+size.x,
			int(w.position.y)+size.y,
		)
		bulletRect := image.Rect(int(b.position.x), int(b.position.y), int(b.position.x)+b.size.x, int(b.position.y)+b.size.y)
		wormRect := image.Rect(int(w.position.x), int(w.position.y), endPosition.x, endPosition.y)
		if isAABBCollision(bulletRect, wormRect) {
			w.health--
		}
	}

}

func (w *Worm) isDead() bool {
	if w.health <= 0 {
		return true
	}
	return false
}

func (w *Worm) kill() {
	// explode
}
