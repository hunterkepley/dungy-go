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
	// UIHealthMeter ... UIIMAGE ENUM [3]
	UIHealthMeter
	// UIEnergyMeter ... UIIMAGE ENUM [4]
	UIEnergyMeter
)

func (u UIImage) String() string {
	return [...]string{
		"Unknown",
		"UIHealthBar",
		"UIEnergyBar",
		"UIHealthMeter",
		"UIEnergyMeter",
	}[u]
}

// ^ UITYPE ENUM ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Generate UI
func generateUI(image *ebiten.Image) []UI {
	return []UI{
		createStaticImage(newVec2i(4, screenHeight-14), UIHealthBar, image),
		createStaticImage(newVec2i(52, screenHeight-14), UIEnergyBar, image),
		createMeterImage(newVec2i(4, screenHeight-14), UIHealthMeter, image),
	}
}

// UI is all the UI in the game
type UI interface {
	render(screen *ebiten.Image)
	uiImage() UIImage
	update(Vec2i)
}

// getUIRect returns a image.Rectangle of the UI image
func getUIRect(ui UIImage) image.Rectangle {
	switch ui {
	case (UIHealthBar):
		return image.Rect(0, 52, 45, 62)
	case (UIEnergyBar):
		return image.Rect(0, 71, 45, 81)
	case (UIHealthMeter):
		return image.Rect(17, 64, 43, 70)
	case (UIEnergyMeter):
		return image.Rect(17, 82, 43, 88)
	}
	// Default, empty
	return image.Rectangle{}
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

func (s StaticImage) update(sizeValues Vec2i) {

}

func (s StaticImage) render(screen *ebiten.Image) {
	subRect := getUIRect(s.ui)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(s.position.x), float64(s.position.y))
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?
	screen.DrawImage(s.image.SubImage(subRect).(*ebiten.Image), op)
}

// MeterImage is a bar, it shrinks or grows based on a value
type MeterImage struct {
	position Vec2i

	ui UIImage

	meterSize Vec2i

	image *ebiten.Image
}

func createMeterImage(position Vec2i, ui UIImage, image *ebiten.Image) *MeterImage {
	return &MeterImage{
		position,
		ui,
		newVec2i(1, 1),
		image,
	}
}

func (m MeterImage) uiImage() UIImage {
	return m.ui
}

// update sizeValues.x = currentSize; sizeValues.y = maxSize
func (m *MeterImage) update(sizeValues Vec2i) {
	m.meterSize = sizeValues
}

func (m MeterImage) render(screen *ebiten.Image) {
	subRect := getUIRect(m.ui)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(m.position.x), float64(m.position.y))
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?
	screen.DrawImage(m.image.SubImage(image.Rect(
		subRect.Min.X,
		subRect.Min.Y,
		subRect.Max.X*(m.meterSize.x/m.meterSize.y),
		subRect.Max.Y,
	)).(*ebiten.Image), op)
}

func updateUI(g *Game) {
	for _, u := range g.ui {
		switch {
		case u.uiImage() == UIHealthMeter:
			u.update(newVec2i(g.player.health, g.player.maxHealth))
		case u.uiImage() == UIEnergyMeter:
			u.update(newVec2i(g.player.energy, g.player.maxEnergy))
		default:
			go u.update(Vec2i{})
		}
	}
}

func renderUI(g *Game, screen *ebiten.Image) {
	for _, u := range g.ui {
		u.render(screen)
	}
}
