package main

import "github.com/hajimehoshi/ebiten"

// Shadow is an image that gets places below an object
type Shadow struct {
	position Vec2f

	sprite Sprite

	image *ebiten.Image
}

func createShadow() Shadow {
	return Shadow{}
}
