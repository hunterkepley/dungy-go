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
	walls  []Tile
	tiles  [][]Tile
}

// Init initializes the game
func (g *Game) Init() {
	// Player
	g.player = createPlayer(newVec2f(screenWidth/2, screenHeight/2))
	// Background image
	testBackgroundImage, _ = loadImage("./Assets/Art/background.png")
	// Fonts
	g.InitFonts()
	// Generate starting walls
	g.walls = generateWalls(itileSpritesheet)
	g.tiles = generateTiles(itileSpritesheet)
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

	for i := 0; i < len(g.walls); i++ {
		g.walls[i].render(screen)
	}
	for i := 0; i < len(g.tiles); i++ {
		for j := 0; j < len(g.tiles[i]); j++ {
			g.tiles[i][j].render(screen)
		}
	}

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
