package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

const (
	screenWidth  = 1366 / 3 // Multiplied by 2 later to scale images
	screenHeight = 768 / 3  // ^
	version      = "0.0.1"
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
	mdataFont, err = loadTTF("./Assets/Font/PressStart2P-Regular.ttf", dpi, 8)
	if err != nil {
		log.Fatal(err)
	}
	mversionFont, err = loadTTF("./Assets/Font/PressStart2P-Regular.ttf", dpi, 8)
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
	return nil
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	bgop := &ebiten.DrawImageOptions{}
	screen.DrawImage(testBackgroundImage, bgop)
	g.player.render(screen)

	// Basic text render calls
	displayVersion(screen)
	displayTPS(screen)
}

// Layout is the screen layout?...
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func displayGameInfo(screen *ebiten.Image) {
	// Draw DUNGY V...
	versionFontPosition := newVec2i(2, 10)
	msg := fmt.Sprintf("DUNGY v%s", version)
	text.Draw(screen, msg, mversionFont, versionFontPosition.x, versionFontPosition.y, color.White)
	// Draw info
	tpsFontPosition := newVec2i(2, 20)
	msg = fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS())
	text.Draw(screen, msg, mdataFont, tpsFontPosition.x, tpsFontPosition.y, color.White)
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
