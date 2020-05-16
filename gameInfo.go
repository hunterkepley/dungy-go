package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
)

var (
	displayInfo          bool = false
	canChangeDisplayInfo bool = true    // If you can press f5 again
	version                   = "0.0.3" // Game version
	title                     = "DUNGY" // Game title
)

func displayGameInfo(screen *ebiten.Image, player Player) {
	// Draw box
	op := &ebiten.DrawImageOptions{}
	informationBoxPosition := newVec2f(0, 0)
	op.GeoM.Translate(informationBoxPosition.x, informationBoxPosition.y)
	screen.DrawImage(iinformationBox, op)
	// Draw DUNGY V...
	versionFontPosition := newVec2i(2, 10)
	msg := fmt.Sprintf("%s v%s", title, version)
	text.Draw(screen, msg, mversionFont, versionFontPosition.x, versionFontPosition.y, color.White)
	// Draw info
	tpsFontPosition := newVec2i(2, 20)
	msg = fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS())
	text.Draw(screen, msg, mdataFont, tpsFontPosition.x, tpsFontPosition.y, color.White)
	// Draw player movement
	movementFontPosition := newVec2i(2, 30)
	msg = fmt.Sprintf("Movement: %s", player.movement)
	text.Draw(screen, msg, mdataFont, movementFontPosition.x, movementFontPosition.y, color.White)

}

// Check if the game should display game info or not
func checkChangeDisplayInfo() {
	if !ebiten.IsKeyPressed(ebiten.KeyF5) {
		canChangeDisplayInfo = true
	} else {
		if canChangeDisplayInfo {
			displayInfo = !displayInfo
			canChangeDisplayInfo = false
		}
	}
}
