package main

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

// PhaseChunk is a chunk of a "map" to phase in randomly
type PhaseChunk struct {
	tiles    []PhaseTile
	portal   Portal
	mapNodes []string
}

func createPhaseChunk(tiles []PhaseTile, portal Portal, mapNodes []string) PhaseChunk {
	return PhaseChunk{
		tiles:    tiles,
		portal:   portal,
		mapNodes: mapNodes,
	}
}

// PhaseTile is a stripped down Tile struct for easier usage with the phase chunk struct
type PhaseTile struct {
	tileType  TileType
	image     *ebiten.Image
	imageRect image.Rectangle

	finishedMoving bool // Whether or not the tile moved down then up like it should
	empty          bool // Whether or not this tile actually changes
}

func createPhaseTile(tileType TileType, image *ebiten.Image, imageRect image.Rectangle, empty bool) PhaseTile {
	return PhaseTile{
		tileType:  tileType,
		image:     image,
		imageRect: imageRect,

		empty: empty,
	}
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

		timer:    200,
		timerMax: 200,
	}

	p.makePhases()

	return p
}

/*
 * This function handles the phase changes, or, the chunks of new maps that replace the current one
 *
 */
func (p *Phases) phaseHandler() {
	maxEnemies := 500
	if p.timer > 0 {
		p.timer--
	} else {
		if len(gameReference.enemies) < maxEnemies {
			p.phase()
		}
		p.timer = p.timerMax
	}
}

func (p *Phases) getRandomPhase() PhaseChunk {
	if len(p.chunks) > 0 {
		return p.chunks[rand.Intn(len(p.chunks))]
	}
	return PhaseChunk{}
}

func (p *Phases) addPhaseChunk(tiles []PhaseTile, portal Portal, mapNodes []string) {
	c := createPhaseChunk(tiles, portal, mapNodes)
	p.chunks = append(p.chunks, c)
}

func (p *Phases) phase() {
	// Get random chunk
	chosenChunk := p.getRandomPhase()

	// Choose random tile and get the index
	_, index := getRandomTile(gameReference, Vec2i{len(chosenChunk.mapNodes[0]) - 1, len(chosenChunk.mapNodes) - 1})

	portalCounter := 0

	currentTile := Vec2i{0, 0}

	if len(chosenChunk.tiles) > 0 {

		for i := 0; i < len(chosenChunk.mapNodes); i++ {
			for j := 0; j < len(chosenChunk.mapNodes[i]); j++ {

				if index.x+i < len(gameReference.tiles) &&
					index.y+j < len(gameReference.tiles[i]) &&
					gameReference.tiles[index.x+j][index.y+i].tileType == SmallTile &&
					i+j < len(chosenChunk.tiles) {

					currentTile = Vec2i{index.x + j, index.y + i}

					gameReference.tiles[currentTile.x][currentTile.y].imageRect = chosenChunk.tiles[i+j].imageRect

					// Portal
					if portalCounter > 2 {
						chosenChunk.portal = createPortal(Vec2f{
							gameReference.tiles[currentTile.x][currentTile.y].position.x,
							gameReference.tiles[currentTile.x][currentTile.y].position.y,
						})

						gameReference.portals = append(gameReference.portals, chosenChunk.portal)
						portalCounter = -100
					}
					portalCounter++
				}

			}
		}

	}
}

// PHASES -- May move to another file? Maybe it's own filetype and parser?

func (p *Phases) makePhases() {
	redTile := image.Rect(119, 0, 135, 17)
	p.addPhaseChunk(
		[]PhaseTile{
			createPhaseTile(SmallTile, itileSpritesheet, redTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, redTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, redTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, redTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, redTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, redTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, redTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, redTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, redTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, redTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, redTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, redTile, true),
		},
		Portal{},
		[]string{
			"    ",
			"    ",
			"    ",
		},
	)
	orangeTile := image.Rect(102, 0, 118, 17)
	p.addPhaseChunk(
		[]PhaseTile{
			createPhaseTile(SmallTile, itileSpritesheet, orangeTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, orangeTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, orangeTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, orangeTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, orangeTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, orangeTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, orangeTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, orangeTile, true),
		},
		Portal{},
		[]string{
			"  ",
			"  ",
			"  ",
			"  ",
		},
	)
}
