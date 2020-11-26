package main

import (
	"github.com/hajimehoshi/ebiten"
)

// Map is a map, it contains MapData, tiles, etc
type Map struct {
	tiles    []*Tile
	lights   []*Light // Background/level-specific lights
	mapNodes []string
	phases   Phases
}

func initMaps(g *Game) {
	gameReference.maps = append(gameReference.maps, Map{})
	initMapSpaceship(g)
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
	gameReference.maps[0].phases.phaseHandler()
}

func (m *Map) render(screen *ebiten.Image) {
}

func initMapSpaceship(g *Game) {
	index := 0
	gameReference.maps[index].phases = initPhases()
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

	gameReference.maps[index].mapNodes = []string{
		"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		"x                            x",
		"x                            x",
		"x                            x",
		"x                            x",
		"x                            x",
		"x                            x",
		"x                            x",
		"x                            x",
		"x                            x",
		"x                            x",
		"x                            x",
		"x                            x",
		"x                            x",
		"x                            x",
		"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
	}

	g.currentMap = gameReference.maps[index]
}
