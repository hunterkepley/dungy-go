package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 1366 / 2 // Multiplied by 2 later to scale images
	screenHeight = 768 / 2  // ^
)

var ()

// Game is the info for the game
type Game struct {
	player Player
}

// Update updates the game
func (g *Game) Update(screen *ebiten.Image) error {
	g.player.update()
	return nil
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	g.player.render(screen)
}

// Layout is the screen layout?...
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func loadPregameResources() {
	loadPlayerImages()
}

func main() {

	loadPregameResources()
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("D U N G Y")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
