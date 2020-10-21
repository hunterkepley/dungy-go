package main

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

// PhaseChunk is a chunk of a "map" to phase in randomly
type PhaseChunk struct {
	tiles    []PhaseTile
	enemies  []Enemy
	mapNodes []string
}

func createPhaseChunk(tiles []PhaseTile, enemies []Enemy, mapNodes []string) PhaseChunk {
	return PhaseChunk{
		tiles:    tiles,
		enemies:  enemies,
		mapNodes: mapNodes,
	}
}

// PhaseTile is a stripped down Tile struct for easier usage with the phase chunk struct
type PhaseTile struct {
	tileType  TileType
	image     *ebiten.Image
	imageRect image.Rectangle
	empty     bool // Whether or not this tile actually changes
}

func createPhaseTile(tileType TileType, image *ebiten.Image, imageRect image.Rectangle, empty bool) PhaseTile {
	return PhaseTile{tileType, image, imageRect, empty}
}

// Phases is a struct containing all the phases in the game and some controls for them
type Phases struct {
	chunks []PhaseChunk

	timer    int
	timerMax int
}

func initPhases() Phases {
	p := Phases{
		chunks: []PhaseChunk{},

		timer:    50,
		timerMax: 50,
	}

	p.makePhases()

	return p
}

/*
 * This function handles the phase changes, or, the chunks of new maps that replace the current one
 *
 * Lots to do
 */
func (p *Phases) phaseHandler() {
	if p.timer > 0 {
		p.timer--
	} else {
		p.phase()
		p.timer = p.timerMax
	}
}

func (p *Phases) makePhases() {
	testTilePosition := image.Rect(102, 0, 118, 17)
	p.addPhaseChunk(
		[]PhaseTile{
			createPhaseTile(SmallTile, itileSpritesheet, testTilePosition, true),
			createPhaseTile(SmallTile, itileSpritesheet, testTilePosition, true),
			createPhaseTile(SmallTile, itileSpritesheet, testTilePosition, true),
			createPhaseTile(SmallTile, itileSpritesheet, testTilePosition, true),
			createPhaseTile(SmallTile, itileSpritesheet, testTilePosition, true),
		},
		[]Enemy{},
		[]string{
			"xxxx",
			"xxxx",
			"xxxx",
			"xxxx",
		},
	)
}

func (p *Phases) getRandomPhase() PhaseChunk {
	if len(p.chunks) > 0 {
		return p.chunks[rand.Intn(len(p.chunks))]
	}
	return PhaseChunk{}
}

func (p *Phases) addPhaseChunk(tiles []PhaseTile, enemies []Enemy, mapNodes []string) {
	c := createPhaseChunk(tiles, enemies, mapNodes)
	p.chunks = append(p.chunks, c)
}

func (p *Phases) phase() {
	// Get random chunk
	chosenChunk := p.getRandomPhase()

	// Choose random tile
	chosenTile, index := getRandomTile(gameReference)

	if len(chosenChunk.tiles) > 0 {
		for i := 0; i < len(chosenChunk.mapNodes); i++ {
			for j := 0; j < len(chosenChunk.mapNodes[i]); j++ {
				if index.x+i < len(gameReference.tiles) &&
					index.y+j < len(gameReference.tiles[i]) &&
					gameReference.tiles[index.x+i][index.y+j].tileType == SmallTile &&
					i+j < len(chosenChunk.tiles) {
					gameReference.tiles[index.x+i][index.y+j].imageRect = chosenChunk.tiles[i+j].imageRect
				}

			}
		}
	}

	_ = chosenTile
	_ = chosenChunk
}
