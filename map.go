package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

// Map is a map, it contains MapData, tiles, etc
type Map struct {
	tiles    [][]Tile
	lights   []*Light // Background/level-specific lights
	mapNodes []string
	phases   Phases
	portals  []Portal
}

func initMaps(g *Game) {
	for i := 0; i < 2; i++ { // Two total maps
		g.maps = append(gameReference.maps, Map{})
	}
	// Spaceship
	g.maps[0] = initMapSpaceship1() // Spaceship room 1
	g.maps[1] = initMapSpaceship2() // Spaceship room 2

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

func switchMap(g *Game, mapNumber int) {
	g.currentMap = g.maps[mapNumber] // Switch current map

	g.currentMap.phases = Phases{}    // Reset phases
	g.currentMap.portals = []Portal{} // Reset portals
	g.enemies = []Enemy{}             // Reset enemies
}

func initMapSpaceship1() Map {
	index := 0 // First map

	gameReference.maps[index].phases = initPhases()

	gameReference.maps[index].tiles = generateTiles(itileSpritesheet) // Tiles

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
		"                              ",
		"                              ",
		"                              ",
		"                              ",
		"                              ",
		"                              ",
		"                              ",
		"                              ",
		"                              ",
		"                              ",
		"                              ",
		"                              ",
		"                              ",
		"                              ",
	}

	return gameReference.maps[index]
}

func initMapSpaceship2() Map {
	index := 1 // Second map

	gameReference.maps[index].phases = initPhases()

	gameReference.maps[index].tiles = generateTiles(itileSpritesheet) // Tiles
	// Lights
	for i := 0; i < 6; i++ {
		offset := 120.
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
		"              xx              ",
		"              xx              ",
		"              xx              ",
		"              xx              ",
		"              xx              ",
		"       xxxxxxxxxxxxxxxx       ",
		"                              ",
		"                              ",
		"              xx              ",
		"              xx              ",
		"            xxxxxx            ",
		"              xx              ",
		"              xx              ",
		"              xx              ",
	}

	gameReference.maps[index].tiles = createWallsFromMap(gameReference.maps[index], itileSpritesheet)

	return gameReference.maps[index]
}
