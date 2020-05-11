package main

import (
	"fmt"
	"image/color"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"

	"github.com/hunterkepley/commonpixel"
)

var (
	// Basic variables
	frames    = 0 // Fps
	second    = time.Tick(time.Second)
	gameState = 1 // 0 is in a game, 1 is in the menu. Keeps track of rendering and updating.
	dt        float64

	imageScale       = 2.
	winWidth         = 1024.
	winHeight        = 768.
	currentWinWidth  = winWidth
	currentWinHeight = winHeight

	//player Player

	viewMatrix pixel.Matrix

	currentShader string

	clearColor color.Color
)

const ()

func run() {
	cfg := pixelgl.WindowConfig{ // Defines window struct
		Title:     "DUNGY",
		Bounds:    pixel.R(0, 0, winWidth, winHeight),
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg) // Creates window
	if err != nil {                    // Deals with error
		panic(err)
	}

	// Set vsync, temporary
	//win.SetVSync(true)

	viewCanvas := pixelgl.NewCanvas(pixel.R(win.Bounds().Min.X, win.Bounds().Min.Y, win.Bounds().W(), win.Bounds().H()))

	loadResources()

	// Set up the matrices for the view of the world
	letterBox(win)

	//player = createPlayer()

	// Set up all levels
	loadLevels()

	viewCanvas.SetFragmentShader(regularShader)

	clearColor = color.RGBA{0x0a, 0x0a, 0x0a, 0x0a}

	last := time.Now()  // For fps decoupled updates
	for !win.Closed() { // Game loop
		if currentWinHeight != win.Bounds().H() || currentWinWidth != win.Bounds().W() {
			// Resize event
			currentWinWidth = win.Bounds().W()
			currentWinHeight = win.Bounds().H()
			letterBox(win)
		}
		imd := imdraw.New(nil)
		dt = time.Since(last).Seconds() // For fps decoupled updates.


		// This is used for when the window is frozen for a very long time to prevent noclip
		if dt > 0.25 { 
			dt = 0.
		}
		last = time.Now() // For fps decoupled updates
		win.Clear(colornames.Black)
		viewCanvas.Clear(clearColor)
		imd.Clear()

		// TODO: For debugging purposes
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			fmt.Println(win.MousePosition())
		}

		switch gameState {
		case 0: // In game, will probably change... Not sure
			updateGame(win, viewCanvas, dt)
			renderGame(win, viewCanvas, imd, dt)
			clearColor = color.RGBA{0x0a, 0x0a, 0x0a, 0x0a}
		case 1: // In menu [?Likely to be separate menus?]
			updateMenu(win, viewCanvas, dt)
			renderMenu(win, viewCanvas)
			clearColor = color.Black
		}

		viewCanvas.Draw(win, viewMatrix)
		imd.Draw(win)
		win.Update()

		frames++ // FPS is dealt here
		select { // Waits for the block to finish
		case <-second: // A second has passed
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames)) // Appends fps to title
			frames = 0                                                   // Reset it my dude
		default:
		} // FPS is done
	}
}

func loadResources() {
	//Load the player sprite sheets for the game
	loadPlayerSpritesheets()
}

func main() {
	pixelgl.Run(run)
}
