package main

import (
	"fmt"
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

// Gib is random piece of a subimage that has a bloodEmitter
type Gib struct {
	position Vec2f
	size     Vec2i
	velocity Vec2i

	bloodEmitter BloodEmitter

	distanceAllowed int
	canMove         bool

	rotation float64

	subImage image.Rectangle
	image    *ebiten.Image
}

func (g *Gib) update(game *Game) {
	notInWall := true
	if g.distanceAllowed >= 0 {
		if g.position.x <= 17+float64(g.size.x) ||
			g.position.y <= 27 ||
			g.position.x+float64(g.size.x) >= screenWidth-17-float64(g.size.x) ||
			g.position.y+float64(g.size.y) >= screenHeight-17-float64(g.size.y) {

			notInWall = false
			g.canMove = false
			if notInWall {
				g.velocity.x *= -1
				g.velocity.y *= -1
			}
		}
		if notInWall {
			// If positive, subtract until 0
			if g.velocity.x != 0 {
				if g.velocity.x > 0 {
					g.velocity.x--
				} else if g.velocity.x < 0 { // Negative
					g.velocity.y++
				}
			}
			if g.velocity.y != 0 {
				if g.velocity.y > 0 {
					g.velocity.y--
				} else if g.velocity.y < 0 { // Negative
					g.velocity.y++
				}
			}
			g.position.x += float64(g.velocity.x)
			g.position.y += float64(g.velocity.y)
			g.distanceAllowed--
		}
	} else {
		g.canMove = false
	}
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

func createGib(position Vec2f,
	rotation float64,
	distanceAllowed int,
	randomVelocity Vec2i,
	subImage image.Rectangle,
	image *ebiten.Image) Gib {

	size := newVec2i(subImage.Max.X-subImage.Min.X, subImage.Max.Y-subImage.Min.Y)
	return Gib{
		position,
		size,
		randomVelocity,
		createBloodEmitter(position, 5, size, ibloodSpritesheet),
		distanceAllowed,
		true,
		rotation,
		subImage,
		image,
	}
}

func (g *GibHandler) update(game *Game) {
	for i := 0; i < len(g.gibs); i++ {
		g.gibs[i].update(game)
		if g.gibs[i].canMove {
			g.gibs[i].bloodEmitter.update(g.gibs[i].position)
		}
	}
}

func (g *GibHandler) render(screen *ebiten.Image) {
	for i := 0; i < len(g.gibs); i++ {
		g.gibs[i].render(screen)
		g.gibs[i].bloodEmitter.render(screen)
	}
}

func (g *GibHandler) explode(numberOfGibs int,
	gibSize int,
	originPosition Vec2f,
	subImage image.Rectangle,
	gibImage *ebiten.Image) {

	for i := 0; i < numberOfGibs; i++ {
		randomDistanceAllowed := 15 + rand.Intn(3)
		randomRotation := float64(rand.Intn(15))
		randomVelocity := newVec2i(rand.Intn(10), rand.Intn(10)) // Random velocity
		switch rand.Intn(4) {
		case (0):
			randomVelocity.x *= -1
		case (1):
			randomVelocity.y *= -1
		case (2):
			randomVelocity.x *= -1
			randomVelocity.y *= -1
		}

		// Get the subimage size and position for random gibs
		subImageSize := newVec2i(subImage.Max.X-subImage.Min.X, subImage.Max.Y-subImage.Min.Y)
		subImagePosition := newVec2i(
			subImage.Min.X+rand.Intn(subImageSize.x-gibSize),
			subImage.Min.Y+rand.Intn(subImageSize.y-gibSize),
		)

		newSubImage := image.Rect(
			subImagePosition.x,
			subImagePosition.y,
			subImagePosition.x+gibSize,
			subImagePosition.y+gibSize,
		)

		g.gibs = append(g.gibs, createGib(
			originPosition,
			randomRotation,
			randomDistanceAllowed,
			randomVelocity,
			newSubImage,
			gibImage,
		))
	}
}

func updateGibHandlers(g *Game) {
	for i := 0; i < len(g.gibHandlers); i++ {
		g.gibHandlers[i].update(g)
		fmt.Println(len(g.gibHandlers) * len(g.gibHandlers[i].gibs))
	}
}

func renderGibHandlers(g *Game, screen *ebiten.Image) {
	for i := 0; i < len(g.gibHandlers); i++ {
		g.gibHandlers[i].render(screen)
	}
}
