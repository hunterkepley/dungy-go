package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 1366 / 3 // Multiplied by 2 later to scale images
	screenHeight = 768 / 3  // ^
)

var (
	gameInitialized     bool = false
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

func loadPregameResources() {
	loadPlayerImages()
}

func main() {

	loadPregameResources()
	ebiten.SetWindowSize(screenWidth*3, screenHeight*3)
	ebiten.SetWindowTitle("D U N G Y")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
