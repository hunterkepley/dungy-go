package main

import (
	"github.com/hajimehoshi/ebiten"
)

// MapData is the metadata for a Map (mainly the name and version)
type MapData struct {
	name    string
	version string
}

// Map is a map, it contains MapData, tiles, etc
type Map struct {
	data MapData

	tiles    []Tile
	lights   []*Light // Background/level-specific lights
	mapNodes []string
}

// I plan for this to end up loading map files just based off map names.
// Most likely will be random like items except an even likelihood of every map
func initMaps(g *Game) {
	gameReference.maps = append(gameReference.maps, Map{data: MapData{"SpaceShip", "0.1"}}) // SpaceShip
	initMapSpaceShip(g)
}

func updateMaps() {
	for i := 0; i < len(gameReference.maps); i++ {
		gameReference.maps[i].update()
	}
}

func renderMaps(screen *ebiten.Image) {
	for i := 0; i < len(gameReference.maps); i++ {
		gameReference.maps[i].render(screen)
	}
}

func (m *Map) update() {
}

func (m *Map) render(screen *ebiten.Image) {
}

// SpaceShip | INDEX 0
func initMapSpaceShip(g *Game) {
	index := 0
	// Lights
	for i := 0; i < 12; i++ {
		offset := 17.
		lightPosition := newVec2f(offset+float64(i)*50, 25)
		lightRotation := 0.
		rect := gameReference.lightHandler.lightImages.rectangleLight1
		id := gameReference.lightHandler.addLightStatic(rect, lightRotation, lightPosition) // Create light
		// Add it to the map's light reference array
		gameReference.maps[index].lights = append(gameReference.maps[index].lights,
			&gameReference.lightHandler.lights[gameReference.lightHandler.getLightIndex(id)],
		)
	}

	// Redo later with a mapNode maker function? Or, manually until editor?
	gameReference.maps[index].mapNodes = []string{
		"xxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		"x                           x",
		"x                           x",
		"x                           x",
		"x                           x",
		"x                           x",
		"x                           x",
		"x                           x",
		"x                           x",
		"xxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
	}
	g.currentMap = gameReference.maps[index]
}
