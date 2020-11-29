package main

import (
	"math/rand"

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
	g.maps = append(gameReference.maps, Map{})
	g.maps[0] = initMapSpaceship() // Spaceship
	g.currentMap = g.maps[0]
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

// Returns a random position that's on a walkable node (' ')
func (m *Map) randomPosition() Vec2f {
	position := Vec2f{0, 0}
	randomTile := Vec2i{0, 0}
	for m.mapNodes[randomTile.x][randomTile.y] != ' ' {
		randomTile = Vec2i{rand.Intn(len(m.mapNodes) - 2), rand.Intn(len(m.mapNodes[0]) - 2)}
	}
	position = Vec2f{float64(randomTile.y * smallTileSize.x), float64(randomTile.x * smallTileSize.y)}

	return position
}

func initMapSpaceship() Map {
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

	return gameReference.maps[index]
}
