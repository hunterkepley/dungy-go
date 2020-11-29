package main

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

// Spawner spawns enemies in on an interval and deals with the spawn visuals
type Spawner struct {
	bits []image.Rectangle

	speed float64

	image *ebiten.Image
}
