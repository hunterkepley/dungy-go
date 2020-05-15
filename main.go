package main

import (
	"fmt"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 1366 / 3 // Multiplied by 2 later to scale images
	screenHeight = 768 / 3  // ^
)

var (
	gameInitialized     bool    = false
	tpsDisplayTimer     float64 = 0
	testBackgroundImage *ebiten.Image
)

// Game is the info for the game
type Game struct {
	player  Player
	player2 Player
}

// Init initializes the game
func (g *Game) Init() {
	g.player = createPlayer(newVec2f(screenWidth/2, screenHeight/2))
	testBackgroundImage, _ = loadImage("./Assets/Art/testBackground.png")
}

// Update updates the game
func (g *Game) Update(screen *ebiten.Image) error {
	if !gameInitialized {
		g.Init()
		gameInitialized = true
	}
	g.player.update()
	displayTPS()
	return nil
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	bgop := &ebiten.DrawImageOptions{}
	screen.DrawImage(testBackgroundImage, bgop)
	g.player.render(screen)
}

// Layout is the screen layout?...
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func displayTPS() {
	// Update TPS
	tpsDisplayTimer++
	if tpsDisplayTimer > 30 {
		ebiten.SetWindowTitle(fmt.Sprint("D U N G Y | FPS: ", ebiten.CurrentTPS()))
		tpsDisplayTimer = 0
	}
}

func loadPregameResources() {
	loadPlayerImages()
}

func main() {

	loadPregameResources()
	ebiten.SetWindowSize(screenWidth*3, screenHeight*3)
	ebiten.SetWindowTitle("D U N G Y")
	ebiten.SetWindowResizable(true)

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
