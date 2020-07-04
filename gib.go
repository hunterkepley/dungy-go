package main

import "github.com/hajimehoshi/ebiten"

type Gib struct {
	position Vec2f

	image *ebiten.Image
}

func createGib() Gib {
	return Gib{}
}

func (g *Gib) update() {

}

func (g *Gib) render(screen *ebiten.Image) {

}
