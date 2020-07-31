package main

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten"
)

// LightImages is a struct of Rectangles of all the light images
type LightImages struct {
	playerLight image.Rectangle
}

// this returns a LightImages struct to init the one in the Game struct
func initLightImages() LightImages {
	return LightImages{
		playerLight: image.Rect(0, 0, 88, 81),
	}
}

// LightBackground is the "background" [the darkness color] of light
type LightBackground struct {
	position Vec2f
	image    *ebiten.Image // the "background" light image
}

func createLightBackground() LightBackground {
	return LightBackground{
		position: newVec2f(0, 0),
		image:    ilightingBackground,
	}
}

func (b *LightBackground) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.position.x, b.position.y)
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?
	screen.DrawImage(b.image, op)
}

func (b *LightBackground) update() {
}

// Light is a circular image with a position that masks with the LightBackground
type Light struct {
	position Vec2f

	id int

	subImage image.Rectangle
	image    *ebiten.Image
}

func createLight(lightRect image.Rectangle, id int) Light {
	return Light{
		subImage: lightRect,
		image:    ilightingSpritesheet,
		id:       id,
	}

}

func (l *Light) update(center Vec2f) {
	size := newVec2i(l.subImage.Max.X-l.subImage.Min.X, l.subImage.Max.Y-l.subImage.Min.Y)
	l.position = newVec2f(center.x-float64(size.x/2), center.y-float64(size.y/2))
}

// LightHandler controls all lights in the game
type LightHandler struct {
	bg LightBackground
	fg LightBackground

	lights []Light

	lightImages LightImages // All light subimages

	lightID int // Keeps track of light ID's, assigns a unique # to each light

	maskedFgImage *ebiten.Image
}

// this returns a lighthandler to init the one in the Game struct
func initLightHandler() LightHandler {
	maskedFg, _ := ebiten.NewImage(screenWidth, screenHeight, ebiten.FilterDefault)
	return LightHandler{
		bg:            createLightBackground(),
		fg:            createLightBackground(),
		lightImages:   initLightImages(),
		lightID:       0,
		maskedFgImage: maskedFg,
	}
}

func (h *LightHandler) render(screen *ebiten.Image) {
	// Reset the bg.
	//h.maskedFgImage.Fill(color.White)
	h.maskedFgImage.Clear()
	op := &ebiten.DrawImageOptions{}
	for i := 0; i < len(h.lights); i++ {
		op.CompositeMode = ebiten.CompositeModeCopy

		op.GeoM.Translate(float64(h.lights[i].position.x), float64(h.lights[i].position.y))
		h.maskedFgImage.DrawImage(h.lights[i].image, op)
	}
	op = &ebiten.DrawImageOptions{}
	op.CompositeMode = ebiten.CompositeModeSourceIn
	h.maskedFgImage.DrawImage(screen, op)

	//screen.Fill(color.RGBA{0x00, 0x00, 0x00, 0xff})
	screen.DrawImage(h.bg.image, &ebiten.DrawImageOptions{})
	screen.DrawImage(h.maskedFgImage, &ebiten.DrawImageOptions{})
}

func (h *LightHandler) update() {

}

func (h *LightHandler) addLight(subImage image.Rectangle) {
	h.lights = append(h.lights, createLight(subImage, h.generateUniqueLightID()))
}

func removeLight(slice []*Light, id int) []*Light { // Removes a light given the ID
	l := -1
	for i := 0; i < len(slice); i++ {
		if slice[i].id == id {
			l = i
		}
	}
	fmt.Print("\nRemoving light with ID ", id)
	return append(slice[:l], slice[l+1:]...)
}

func (h *LightHandler) generateUniqueLightID() int { // Generates a new ID for a light
	h.lightID++
	return h.lightID
}
