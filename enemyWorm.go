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
	position  Vec2f
	center    Vec2f
	size      Vec2i
	velocity  Vec2f
	moveSpeed float64

	health    int
	maxHealth int
	dead      bool
	remove    bool // Do we remove this worm?
	flipped   bool // Is the worm flipped?

	shadow *Shadow // The shadow below the worm

	subImageRect image.Rectangle

	spritesheet     Spritesheet
	animation       Animation
	animations      WormAnimations
	animationSpeeds WormAnimationSpeeds

	image *ebiten.Image
}

func createWorm(position Vec2f, game *Game) *Worm {
	idleFrontSpritesheet := createSpritesheet(newVec2i(0, 0), newVec2i(68, 22), 4, ienemiesSpritesheet)

	shadowRect := image.Rect(0, 231, 14, 237)
	shadow := createShadow(shadowRect, iplayerSpritesheet, generateUniqueShadowID(game))
	game.shadows = append(game.shadows, &shadow)

	return &Worm{
		position:  position,
		velocity:  newVec2f(0, 0),
		moveSpeed: 0.5,

		health:    3,
		maxHealth: 3,
		dead:      false,

		shadow: &shadow,

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
	// FLIP DECIDER
	flip := newVec2f(1, 1)
	if w.flipped {
		flip.x = -1
	}

	// ROTATE & FLIP
	op.GeoM.Translate(float64(0-w.size.x)/2, float64(0-w.size.y)/2)
	op.GeoM.Scale(flip.x, flip.y)
	op.GeoM.Translate(float64(w.size.x)/2, float64(w.size.y)/2)
	w.subImageRect = image.Rect(
		w.spritesheet.sprites[w.animation.currentFrame].startPosition.x,
		w.spritesheet.sprites[w.animation.currentFrame].startPosition.y,
		w.spritesheet.sprites[w.animation.currentFrame].endPosition.x,
		w.spritesheet.sprites[w.animation.currentFrame].endPosition.y,
	)
	// POSITION
	op.GeoM.Translate(float64(w.position.x), float64(w.position.y))
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?

	screen.DrawImage(w.image.SubImage(w.subImageRect).(*ebiten.Image), op) // Draw worm
}

func (w *Worm) update(game *Game) {
	// Start the animation if it's not playing
	if w.animation.state != AnimationPlayingForwards {
		w.animation = w.animations.idleFront
		w.animation.startForwards()
	}
	w.animation.update(w.animationSpeeds.idle)

	// Pathfind to player
	w.followPlayer(game)

	// Move worm
	w.position.x += w.velocity.x
	w.position.y += w.velocity.y

	w.size = newVec2i(
		w.spritesheet.sprites[w.animation.currentFrame].size.x,
		w.spritesheet.sprites[w.animation.currentFrame].size.y,
	)
	endPosition := newVec2i(
		int(w.position.x)+w.size.x,
		int(w.position.y)+w.size.y,
	)

	// Update shadow
	w.shadow.update(w.position, w.size)

	w.subImageRect = image.Rect(int(w.position.x), int(w.position.y), endPosition.x, endPosition.y)
	w.center = newVec2f(w.position.x+float64(w.size.x)/2, w.position.y+float64(w.size.y)/2)
}

func (w *Worm) followPlayer(game *Game) {
	if w.position.x < game.player.position.x-w.moveSpeed {
		w.velocity.x = w.moveSpeed
		w.flipped = true
	} else if w.position.x > game.player.position.x+w.moveSpeed {
		w.velocity.x = -w.moveSpeed
		w.flipped = false
	} else {
		w.velocity.x = 0
	}

	if w.position.y < game.player.position.y-w.moveSpeed {
		w.velocity.y = w.moveSpeed
	} else if w.position.y > game.player.position.y+w.moveSpeed {
		w.velocity.y = -w.moveSpeed
	} else {
		w.velocity.y = 0
	}
}

func (w *Worm) isDead() bool {
	if w.health <= 0 {
		if !w.dead {
			w.dead = true
		}
		return true
	}
	return false
}

func (w *Worm) getCenter() Vec2f {
	return w.center
}

func (w *Worm) getCurrentSubImageRect() image.Rectangle {
	return w.subImageRect
}

func (w *Worm) getImage() *ebiten.Image {
	return w.image
}

func (w *Worm) damage(value int) {
	w.health--
}

func (w *Worm) getShadow() Shadow {
	return *w.shadow
}
