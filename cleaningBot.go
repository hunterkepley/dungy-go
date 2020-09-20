package main

import "github.com/hajimehoshi/ebiten"

// CleaningBot is a robot that cleans blood!
type CleaningBot struct {
	position Vec2f

	rotation float64

	image *ebiten.Image
}
