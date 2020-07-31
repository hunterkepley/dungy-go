package main

import "github.com/hajimehoshi/ebiten"

// LightBackground is the "background" [the darkness color] of light
type LightBackground struct {
	image *ebiten.Image // the "background" light image
}

// Light is a circular image with a position that masks with the LightBackground
type Light struct {
	image *ebiten.Image
}

func createLightBackground() LightBackground {
	return LightBackground{}
}

func createLight() Light {
	return Light{}
}

func (b *LightBackground) render(screen *ebiten.Image) {

}
