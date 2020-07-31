package main

import "github.com/hajimehoshi/ebiten"

// LightBackground is the "background" [the darkness color] of light
type LightBackground struct {
	image *ebiten.Image // the "background" light image
}

func createLightBackground() LightBackground {
	return LightBackground{}
}

func (b *LightBackground) render(screen *ebiten.Image) {

}

func (b *LightBackground) update() {
}

// Light is a circular image with a position that masks with the LightBackground
type Light struct {
	image *ebiten.Image
}

func createLight() Light {
	return Light{}

}

func (l *Light) render(screen *ebiten.Image) {

}

func (l *Light) update() {

}

// LightHandler controls all lights in the game
type LightHandler struct {
	bg     LightBackground
	lights []Light
}

// this returns a lighthandler to init the one in the Game struct
func initLightHandler() LightHandler {
	return LightHandler{}
}

func (h *LightHandler) render(screen *ebiten.Image) {

}

func (h *LightHandler) update() {

}
