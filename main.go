package main

import (
	"image"
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
	cursor Cursor

	walls   []Tile
	tiles   [][]Tile
	borders []Border
	ui      []UI
}

// Init initializes the game
func (g *Game) Init() {
	// Player
	g.player = createPlayer(newVec2f(screenWidth/2, screenHeight/2))
	// Cursor
	g.cursor = createCursor(iUISpritesheet)
	// Background image
	testBackgroundImage, _ = loadImage("./Assets/Art/background.png")
	// Fonts
	g.InitFonts()
	// Generate starting walls
	g.walls = generateWalls(itileSpritesheet)
	g.tiles = generateTiles(itileSpritesheet)
	generateBigTiles(g.tiles, itileSpritesheet)
	g.borders = generateBorders(itileSpritesheet)
	g.ui = generateUI(iUISpritesheet)
}

// Update updates the game
func (g *Game) Update(screen *ebiten.Image) error {
	if !gameInitialized {
		g.Init()
		gameInitialized = true
	}

	// Update cursor
	g.cursor.update()

	// Update player
	g.player.update()

	// Game info update/check
	go checkChangeDisplayInfo()

	// Update UI
	updateUI(g)

	// Temporary
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(true)
	}
	return nil
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	bgop := &ebiten.DrawImageOptions{}
	screen.DrawImage(testBackgroundImage, bgop)

	// Render game walls/tiles
	renderTiles(g, screen)
	// Render borders
	for _, b := range g.borders {
		b.render(screen)
	}

	// Render player
	g.player.render(screen)

	// Basic text render calls
	if displayInfo {
		displayGameInfo(screen, g.player)
	}

	// Render UI
	renderUI(g, screen)

	// Render cursor
	g.cursor.render(screen)
}

// Layout is the screen layout?...
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {

	loadPregameResources()

	ebiten.SetWindowSize(screenWidth*3, screenHeight*3)
	ebiten.SetWindowTitle("U N R A Y")
	ebiten.SetWindowResizable(true)

	// Load icon
	icon16x16, _ := loadRegularImage("./Assets/Art/Icon/icon16.png")
	icon32x32, _ := loadRegularImage("./Assets/Art/Icon/icon32.png")
	icon48x48, _ := loadRegularImage("./Assets/Art/Icon/icon48.png")
	icon64x64, _ := loadRegularImage("./Assets/Art/Icon/icon64.png")
	ebiten.SetWindowIcon([]image.Image{icon16x16, icon32x32, icon48x48, icon64x64})

	// Hide cursor
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	//ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func loadPregameResources() {
	loadPlayerImages()
	loadUIImages()
	loadTileImages()
	loadItemsImages()
}
