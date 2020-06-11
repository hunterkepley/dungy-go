package main

import (
	"github.com/hajimehoshi/ebiten"
)

var (
	idAnimation = 0
)

// AnimationState is a type for the current state of an animation
type AnimationState int

const (
	// AnimationPlayingForwards ... ANIMATIONSTATE ENUM [1]
	AnimationPlayingForwards AnimationState = iota + 1
	// AnimationStopped ... ANIMATIONSTATE ENUM [2]
	AnimationStopped
	// AnimationPlayingBackwards ... ANIMATIONSTATE ENUM [3]
	AnimationPlayingBackwards
)

func (a AnimationState) String() string {
	return [...]string{"Unknown", "AnimationPlayingForwards", "AnimationStopped", "AnimationPlayingBackwards"}[a]
}

// ^ MOVEMENT ENUM ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Animation a system for animating an image using Ebiten subsets
type Animation struct {
	spritesheet  Spritesheet
	currentFrame int
	timer        float64
	maxTimer     float64
	id           int            // For checking if an animation equals another
	state        AnimationState // If the animation is currently playing
}

func createAnimation(spritesheet Spritesheet, image *ebiten.Image) Animation {

	idAnimation++
	return Animation{
		spritesheet,
		0,           // currentFrame
		0,           // timer
		10,          // maxTimer
		idAnimation, // Animation ID
		AnimationPlayingForwards,
	}
}

func (a *Animation) startForwards() {
	a.state = AnimationPlayingForwards
}

func (a *Animation) startBackwards() {
	a.state = AnimationPlayingBackwards
	a.currentFrame = a.spritesheet.numberOfSprites - 1 // Set to end!
}

func (a *Animation) update(speed float64) {
	if a.timer >= a.maxTimer {
		a.timer = 0
		switch a.state {
		case AnimationPlayingForwards: // Forwards
			a.currentFrame++
			if a.currentFrame == a.spritesheet.numberOfSprites {
				// Reached the end!
				a.currentFrame = 0
			}

		case AnimationPlayingBackwards: // Backwards
			a.currentFrame--
			if a.currentFrame == -1 {
				a.currentFrame = a.spritesheet.numberOfSprites - 1 // Set to end again!
			}
		}
	} else {
		a.timer += 1 * speed
	}
}

func (a *Animation) pause() {
	a.state = AnimationStopped
}
