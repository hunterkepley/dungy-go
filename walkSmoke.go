package main

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

// WalkSmokeEmitter is an emitter that emits WalkSmoke
type WalkSmokeEmitter struct {
	spawnRate    float64
	spawnRateMax float64

	particles []WalkSmoke

	image *ebiten.Image
}

func createWalkSmokeEmitter(spawnRate float64, image *ebiten.Image) WalkSmokeEmitter {
	return WalkSmokeEmitter{
		spawnRate:    spawnRate,
		spawnRateMax: spawnRate,

		image: image,
	}
}

func (e *WalkSmokeEmitter) update(position Vec2f, size Vec2i, direction Direction, isMoving bool) {
	if e.spawnRate > 0 {
		e.spawnRate--
	} else {
		if isMoving {
			e.particles = append(e.particles, createWalkSmoke(position, size, direction, e.image))
			e.spawnRate = e.spawnRateMax
		}
	}

	for i := 0; i < len(e.particles); i++ {
		if e.particles[i].animation.currentFrame == e.particles[i].animation.spritesheet.numberOfSprites-1 &&
			e.particles[i].animationCanFinish {
			removeWalkSmoke(e.particles, i)
			continue
		}
		e.particles[i].update()
	}
}

func (e *WalkSmokeEmitter) render(screen *ebiten.Image) {
	for i := 0; i < len(e.particles); i++ {
		e.particles[i].render(screen)
	}
}

// WalkSmoke the smoke that kicks up behind feet
type WalkSmoke struct {
	position Vec2f

	velocity  Vec2f
	moveSpeed Vec2f

	animation          Animation
	animationSpeed     float64
	animationCanFinish bool

	image *ebiten.Image
}

func createWalkSmoke(position Vec2f, size Vec2i, direction Direction, image *ebiten.Image) WalkSmoke {

	moveSpeed := newVec2f(rand.Float64(), -1*rand.Float64())

	startingPosition := newVec2f(position.x, position.y+float64(size.y))

	velocity := newVec2f(moveSpeed.x, moveSpeed.y)

	// Decides velocity based on movement direction
	switch {
	case direction == Left || direction == DownLeft || direction == UpLeft:
		startingPosition.x += float64(size.x)
	case direction == Right || direction == DownRight || direction == UpRight:
		velocity.x *= -1
	case direction == Up:
		startingPosition.x += float64(size.x / 2)
		velocity.x = rand.Float64()/2 - 0.25
		velocity.y *= -1
	case direction == Down:
		startingPosition.x += float64(size.x / 2)
		velocity.x = rand.Float64()/2 - 0.25
	}

	// Offsets to the proper position under the player
	startingPosition.x -= 2
	startingPosition.y -= 4

	spritesheet := createSpritesheet(newVec2i(0, 0), newVec2i(20, 5), 4, image)
	animation := createAnimation(spritesheet, image)

	animationSpeed := 1.2

	return WalkSmoke{
		position: startingPosition,

		animation:          animation,
		animationSpeed:     animationSpeed,
		animationCanFinish: false,

		velocity: velocity,
		image:    image,
	}
}

func (w *WalkSmoke) update() {
	w.animation.update(w.animationSpeed)

	w.position.x += w.velocity.x
	w.position.y += w.velocity.y

	if w.animation.currentFrame == 1 {
		w.animationCanFinish = true
	}
}

func (w *WalkSmoke) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(w.position.x, w.position.y)
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?

	currentFrame := w.animation.spritesheet.sprites[w.animation.currentFrame]

	subImageRect := image.Rect(
		currentFrame.startPosition.x,
		currentFrame.startPosition.y,
		currentFrame.endPosition.x,
		currentFrame.endPosition.y,
	)

	screen.DrawImage(w.image.SubImage(subImageRect).(*ebiten.Image), op)

}

func removeWalkSmoke(slice []WalkSmoke, w int) []WalkSmoke {
	return append(slice[:w], slice[w+1:]...)
}
