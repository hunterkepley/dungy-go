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
	dead      bool

	subImageRect image.Rectangle

	spritesheet     Spritesheet
	animation       Animation
	animations      WormAnimations
	animationSpeeds WormAnimationSpeeds

	gibHandler GibHandler

	image *ebiten.Image
}

func createWorm(position Vec2f) *Worm {
	idleFrontSpritesheet := createSpritesheet(newVec2i(0, 0), newVec2i(68, 22), 4, ienemiesSpritesheet)

	return &Worm{
		position: position,

		health:    10,
		maxHealth: 10,
		dead:      false,

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
	w.subImageRect = image.Rect(
		w.spritesheet.sprites[w.animation.currentFrame].startPosition.x,
		w.spritesheet.sprites[w.animation.currentFrame].startPosition.y,
		w.spritesheet.sprites[w.animation.currentFrame].endPosition.x,
		w.spritesheet.sprites[w.animation.currentFrame].endPosition.y,
	)
	// Draw gibs
	w.gibHandler.render(screen)

	screen.DrawImage(w.image.SubImage(w.subImageRect).(*ebiten.Image), op) // Draw worm
}

func (w *Worm) update(bullets []Bullet) {
	// Start the animation if it's not playing
	if w.animation.state != AnimationPlayingForwards {
		w.animation = w.animations.idleFront
		w.animation.startForwards()
	}
	w.animation.update(w.animationSpeeds.idle)

	w.gibHandler.update()

	//
	for _, b := range bullets {
		size := newVec2i(
			w.spritesheet.sprites[w.animation.currentFrame].size.x,
			w.spritesheet.sprites[w.animation.currentFrame].size.y,
		)
		endPosition := newVec2i(
			int(w.position.x)+size.x,
			int(w.position.y)+size.y,
		)
		wormRect := image.Rect(int(w.position.x), int(w.position.y), endPosition.x, endPosition.y)
		bulletRect := image.Rect(int(b.position.x), int(b.position.y), int(b.position.x)+b.size.x, int(b.position.y)+b.size.y)
		if isAABBCollision(bulletRect, wormRect) {
			w.health--
		}
	}

}

func (w *Worm) isDead() bool {
	if w.health <= 0 {
		if !w.dead {
			w.dead = true
			w.kill()
		}
		return true
	}
	return false
}

func (w *Worm) kill() { // TODO: Separate gibhandlers from enemies, make them just create one upon death
	w.gibHandler.explode(10, w.subImageRect, w.image)
}
