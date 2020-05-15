package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/font"
)

const (
	screenWidth  = 1366 / 3 // Multiplied by 2 later to scale images
	screenHeight = 768 / 3  // ^
)

var (
	gameInitialized     bool = false
	testBackgroundImage *ebiten.Image

	mdataFont    font.Face
	mversionFont font.Face
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
	g.InitFonts()
}

// InitFonts initializes fonts for the game
func (g *Game) InitFonts() {

	const dpi = 72
	var err error
	mdataFont, err = loadTTF("./Assets/Font/LoRe.ttf", dpi, 8)
	if err != nil {
		log.Fatal(err)
	}
	mversionFont, err = loadTTF("./Assets/Font/LoRe.ttf", dpi, 8)
	if err != nil {
		log.Fatal(err)
	}
}

// Update updates the game
func (g *Game) Update(screen *ebiten.Image) error {
	if !gameInitialized {
		g.Init()
		gameInitialized = true
	}
	g.player.update()
	checkChangeDisplayInfo()
	return nil
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	bgop := &ebiten.DrawImageOptions{}
	screen.DrawImage(testBackgroundImage, bgop)
	g.player.render(screen)

	// Basic text render calls
	if displayInfo {
		displayGameInfo(screen)
	}
}

// Layout is the screen layout?...
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func loadPregameResources() {
	loadPlayerImages()
	loadUIImages()
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
