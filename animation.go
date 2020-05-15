package main

import (
	"github.com/hajimehoshi/ebiten"
)

var (
	idAnimation = 0
)

// Animation a system for animating an image using Ebiten subsets
type Animation struct {
	spritesheet  Spritesheet
	currentFrame int
	timer        float64
	maxTimer     float64
	id           int // For checking if an animation equals another
}

func createAnimation(spritesheet Spritesheet, image *ebiten.Image) Animation {

	idAnimation++
	return Animation{
		spritesheet,
		0,           // currentFrame
		0,           // timer
		10,          // maxTimer
		idAnimation, // Animation ID
	}
}

func (a *Animation) play(speed float64) {
	if a.timer >= a.maxTimer {
		a.timer = 0
		a.currentFrame++
		if a.currentFrame == a.spritesheet.numberOfSprites {
			// Reached the end!
			a.currentFrame = 0
		}
	} else {
		a.timer += 1 * speed
	}
}
