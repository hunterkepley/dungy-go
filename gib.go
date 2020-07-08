package main

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

// Gib is random piece of a subimage that has a bloodEmitter
type Gib struct {
	position Vec2f

	rotation float64

	subImage image.Rectangle
	image    *ebiten.Image
}

func (g *Gib) update() {

}

func (g *Gib) render(screen *ebiten.Image) {

}

// GibHandler controls gibs, it's like a particle mitter but for a set of gibs
type GibHandler struct {
	gibs []Gib
}

func createGibHandler(numberOfGibs int, subImage image.Rectangle, image *ebiten.Image) GibHandler {
	var gibs []Gib

	randomPosition := newVec2f(float64(rand.Intn(10)), float64(rand.Intn(10)))
	randomRotation := float64(rand.Intn(100))

	for i := 0; i < numberOfGibs; i++ {
		gibs = append(gibs, createGib(randomPosition, randomRotation, subImage, image))
	}

	return GibHandler{
		gibs,
	}
}

func createGib(position Vec2f, rotation float64, subImage image.Rectangle, image *ebiten.Image) Gib {
	return Gib{
		position,
		rotation,
		subImage,
		image,
	}
}

func (g *GibHandler) update() {
	for _, gib := range g.gibs {
		gib.update()
	}
}

func (g *GibHandler) render(screen *ebiten.Image) {
	for _, gib := range g.gibs {
		gib.render(screen)
	}
}
