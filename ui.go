package main

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

// UIImage is a type for the UIImage enum
type UIImage int

const (
	// UIHealthBar ... UIIMAGE ENUM [1]
	UIHealthBar UIImage = iota + 1
	// UIEnergyBar ... UIIMAGE ENUM [2]
	UIEnergyBar
)

func (u UIImage) String() string {
	return [...]string{"Unknown", "UIHealthBar", "UIEnergyBar"}[u]
}

// ^ UITYPE ENUM ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Generate UI
func generateUI(image *ebiten.Image) []UI {
	return []UI{
		createStaticImage(newVec2i(4, screenHeight-14), UIHealthBar, image),
		createStaticImage(newVec2i(52, screenHeight-14), UIEnergyBar, image),
	}
}

// UI is all the UI in the game
type UI interface {
	render(screen *ebiten.Image)
	uiImage() UIImage
	update()
}

// StaticImage is an image with nothing else special
type StaticImage struct {
	position Vec2i

	ui UIImage

	image *ebiten.Image
}

func createStaticImage(position Vec2i, ui UIImage, image *ebiten.Image) StaticImage {
	return StaticImage{
		position,
		ui,
		image,
	}
}

func (s StaticImage) uiImage() UIImage {
	return s.ui
}

func (s StaticImage) update() {

}

func (s StaticImage) render(screen *ebiten.Image) {
	var subRect image.Rectangle
	switch s.ui {
	case (UIHealthBar):
		subRect = image.Rect(0, 52, 45, 62)
	case (UIEnergyBar):
		subRect = image.Rect(0, 71, 45, 81)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(s.position.x), float64(s.position.y))
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?
	screen.DrawImage(s.image.SubImage(subRect).(*ebiten.Image), op)
}
