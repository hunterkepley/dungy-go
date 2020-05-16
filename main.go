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
	player Player
	t1     Tile
	t2     Tile
	t3     Tile
}

// Init initializes the game
func (g *Game) Init() {
	g.player = createPlayer(newVec2f(screenWidth/2, screenHeight/2))
	testBackgroundImage, _ = loadImage("./Assets/Art/background.png")
	g.InitFonts()
	g.t1 = createTile(newVec2f(0, 30), SmallTile)
	g.t2 = createTile(newVec2f(30, 30), BigTile)
	g.t3 = createTile(newVec2f(70, 30), WallTile)
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

	g.t1.render(screen)
	g.t2.render(screen)
	g.t3.render(screen)

	g.player.render(screen)

	// Basic text render calls
	if displayInfo {
		displayGameInfo(screen, g.player)
	}
}

// Layout is the screen layout?...
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
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

func loadPregameResources() {
	loadPlayerImages()
	loadUIImages()
	loadTileImages()
}
