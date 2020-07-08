package main

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

// Gib is random piece of a subimage that has a bloodEmitter
type Gib struct {
	position Vec2f
	size     Vec2i

	rotation float64

	subImage image.Rectangle
	image    *ebiten.Image
}

func (g *Gib) update() {
}

func (g *Gib) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0-float64(g.size.x)/2, 0-float64(g.size.y)/2)
	op.GeoM.Rotate(g.rotation)
	op.GeoM.Translate(g.position.x, g.position.y)
	screen.DrawImage(g.image.SubImage(g.subImage).(*ebiten.Image), op) // Draw gib
}

// GibHandler controls gibs, it's like a particle mitter but for a set of gibs
type GibHandler struct {
	gibs []Gib
}

func createGibHandler() GibHandler {
	var gibs []Gib

	return GibHandler{
		gibs,
	}
}

func createGib(position Vec2f, rotation float64, subImage image.Rectangle, image *ebiten.Image) Gib {
	size := newVec2i(subImage.Max.X-subImage.Min.X, subImage.Max.Y-subImage.Min.Y)
	return Gib{
		position,
		size,
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

func (g *GibHandler) explode(numberOfGibs int, subImage image.Rectangle, image *ebiten.Image) {
	randomPosition := newVec2f(float64(rand.Intn(100)), float64(rand.Intn(100)))
	randomRotation := float64(rand.Intn(100))

	for i := 0; i < numberOfGibs; i++ {
		g.gibs = append(g.gibs, createGib(randomPosition, randomRotation, subImage, image))
	}
}
