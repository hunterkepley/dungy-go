package main

import (
	"image"
	_ "image/png"
	"log"
	"math/rand"

	paths "github.com/SolarLune/paths"
	"github.com/hajimehoshi/ebiten"
	lua "github.com/yuin/gopher-lua"
	"golang.org/x/image/font"
)

const (
	screenWidth  = 1366 / 3 // Multiplied by 3 later to scale images
	screenHeight = 768 / 3  // ^
)

var (
	gameInitialized     bool = false
	testBackgroundImage *ebiten.Image

	mdataFont    font.Face
	mversionFont font.Face

	gameReference *Game

	tempSpawnCount int = 1
)

// Game is the info for the game
type Game struct {
	player Player
	cursor Cursor

	enemies          []Enemy
	items            []Item
	gibHandlers      []GibHandler
	bloodEmitters    []BloodEmitter
	bulletExplosions []BulletExplosion
	shadows          []*Shadow
	walls            []Tile
	tiles            [][]Tile
	borders          []Border
	ui               []UI
	maps             []Map

	currentMap Map

	lightHandler LightHandler // Controls the lights!

	shadowID int // Shadow IDs, starts at 0 then increments when a shadow is added

	state int // The game state, 0 is in main menu, 1 is in game, 2 is paused

	settings Settings // Game settings
}

// Init initializes the game
func (g *Game) Init() {

	// Set gameReference
	gameReference = g

	// Init lightHandler
	g.lightHandler = initLightHandler()
	playerLightID := g.lightHandler.addLight(g.lightHandler.lightImages.playerLight, 0)

	// Init maps
	initMaps(g)

	// Player
	g.player = createPlayer(
		newVec2f(screenWidth/2, screenHeight/2),
		g,
		playerLightID,
	)
	g.shadows = append(g.shadows, &g.player.shadow)

	// Test items! ============================
	// TODO: REMOVE THIS                      v
	testItem := createItem(
		newVec2f(300, 100), // Position
		iitemsSpritesheet,  // Image
	)
	testItem.init()
	g.items = append(g.items, testItem)
	testItem = createItem(
		newVec2f(200, 100), // Position
		iitemsSpritesheet,  // Image
	)
	testItem.init()
	g.items = append(g.items, testItem)
	testItem = createItem(
		newVec2f(100, 100), // Position
		iitemsSpritesheet,  // Image
	)
	testItem.init()
	g.items = append(g.items, testItem)
	// TODO: REMOVE THIS                      ^
	// Test items! ============================

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

	// Init lua functions
	L := lua.NewState()
	defer L.Close()
	initLuaFunctions(L)

	// Make the astar path channel
	astarChannel = make(chan *paths.Path, 500)

	// Obviously, temporary
	g.enemies = append(g.enemies, Enemy(createBeefEye(newVec2f(float64(rand.Intn(screenWidth)), float64(rand.Intn(screenHeight))), g)))
	//g.enemies = append(g.enemies, Enemy(createBeefEye(newVec2f(float64(rand.Intn(screenWidth)), float64(rand.Intn(screenHeight))), g)))

	// State starts in game [temporary]
	g.state = 1

	// Init music
	loadMusic()
	// Play song
	go music[0].play()

	// GAME SETTINGS
	loadSettings(&g.settings)

	if g.settings.Graphics.Fullscreen { // Enable fullscreen if enabled
		ebiten.SetFullscreen(true)
	}

}

var e = 0

// Update updates the game
func (g *Game) Update(screen *ebiten.Image) error {
	if !gameInitialized {
		g.Init()
		gameInitialized = true
	}

	// Update game
	if g.state == 1 {
		updateGame(screen, g)
	}

	return nil
}

func updateGame(screen *ebiten.Image, g *Game) {
	// Update cursor
	g.cursor.update()

	// Update map
	updateMaps()

	// Update player
	g.player.update(g)
	g.player.gun.updateBullets(g)
	updateBulletExplosions(g)

	// Update items
	updateItems(g)

	// Update enemies
	updateEnemies(g)

	// Update gib handlers
	updateGibHandlers(g)

	// Update light
	g.lightHandler.update()

	// Game info update/check
	checkChangeDisplayInfo()

	// Update UI
	updateUI(g)
}

// Draw renders everything!
func (g *Game) Draw(screen *ebiten.Image) {
	bgop := &ebiten.DrawImageOptions{}
	screen.DrawImage(testBackgroundImage, bgop)

	// Draw game
	if g.state == 1 { // inGame
		drawGame(screen, g)
	}

}

func drawGame(screen *ebiten.Image, g *Game) {
	// Render game walls/tiles
	renderTiles(g, screen)
	// Render maps
	renderMaps(screen)

	// Render gibHandlers
	renderGibHandlers(g, screen)

	renderItems(g, screen)

	// Render shadows!
	for i := 0; i < len(g.shadows); i++ {
		g.shadows[i].render(screen)
	}

	// Render enemies behind player
	renderEnemies(g, screen)

	// Render player
	g.player.render(screen)
	g.player.gun.renderBullets(screen)
	if g.player.isDrawable {
		g.player.gun.render(screen)
	}
	renderBulletExplosions(g, screen)

	// Render borders
	for _, b := range g.borders {
		b.render(screen)
	}

	// Render light
	g.lightHandler.render(screen)

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

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func loadPregameResources() {
	// Images
	loadPlayerImages()
	loadUIImages()
	loadTileImages()
	loadItemsImages()
	loadEnemiesImages()
	loadParticlesImages()
	loadLightingImages()

	initListOfAllItemNames()
}
