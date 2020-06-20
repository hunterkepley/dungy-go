package main

import (
	"github.com/hajimehoshi/ebiten"
)

// Gun is the player's gun
type Gun struct {
	position Vec2f
	size     Vec2i

	image *ebiten.Image
}

func (g *Gun) render(screen *ebiten.Image) {

}

func (g *Gun) update() {

}
