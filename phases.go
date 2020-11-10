package main

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

// PhaseChunk is a chunk of a "map" to phase in randomly
type PhaseChunk struct {
	tiles    []PhaseTile
	enemies  []PhaseEnemy
	mapNodes []string
}

func createPhaseChunk(tiles []PhaseTile, enemies []PhaseEnemy, mapNodes []string) PhaseChunk {
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

	finishedMoving bool // Whether or not the tile moved down then up like it should
	empty          bool // Whether or not this tile actually changes
}

// PhaseEnemy is a stripped down Enemy struct for easier usage with the phase chunk struct
type PhaseEnemy struct {
	enemyType    EnemyType
	tilePosition Vec2i
}

func createPhaseTile(tileType TileType, image *ebiten.Image, imageRect image.Rectangle, empty bool) PhaseTile {
	return PhaseTile{
		tileType:  tileType,
		image:     image,
		imageRect: imageRect,

		finishedMoving: false,
		empty:          empty,
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

		timer:    1,
		timerMax: 1,
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
		if len(gameReference.enemies) < 500 {
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

func (p *Phases) addPhaseChunk(tiles []PhaseTile, enemies []PhaseEnemy, mapNodes []string) {
	c := createPhaseChunk(tiles, enemies, mapNodes)
	p.chunks = append(p.chunks, c)
}

func (p *Phases) phase() {
	// Get random chunk
	chosenChunk := p.getRandomPhase()

	// Choose random tile and get the index
	_, index := getRandomTile(gameReference, Vec2i{len(chosenChunk.mapNodes[0]), len(chosenChunk.mapNodes[1])})

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

					chosenChunk.tiles[i+j].finishedMoving = true
					// Enemies
					for e := 0; e < len(chosenChunk.enemies); e++ {
						if i == chosenChunk.enemies[e].tilePosition.x &&
							j == chosenChunk.enemies[e].tilePosition.y {

							generateEnemy(
								chosenChunk.enemies[e].enemyType,
								newVec2f(gameReference.tiles[currentTile.x][currentTile.y].position.x, gameReference.tiles[currentTile.x][currentTile.y].position.y),
								gameReference,
							)
						}

					}
				}

			}
		}

	}
}

// PHASES -- May move to another file? Maybe it's own filetype and parser?

func (p *Phases) makePhases() {
	beefEyeTile := image.Rect(119, 0, 135, 17)
	p.addPhaseChunk(
		[]PhaseTile{
			createPhaseTile(SmallTile, itileSpritesheet, beefEyeTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, beefEyeTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, beefEyeTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, beefEyeTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, beefEyeTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, beefEyeTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, beefEyeTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, beefEyeTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, beefEyeTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, beefEyeTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, beefEyeTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, beefEyeTile, true),
		},
		[]PhaseEnemy{
			{
				EBeefEye,
				Vec2i{0, 0}, // Spawn on first tile
			},
			{
				EBeefEye,
				Vec2i{3, 3},
			},
		},
		[]string{
			"    ",
			"    ",
			"    ",
		},
	)
	wormTile := image.Rect(102, 0, 118, 17)
	p.addPhaseChunk(
		[]PhaseTile{
			createPhaseTile(SmallTile, itileSpritesheet, wormTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, wormTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, wormTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, wormTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, wormTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, wormTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, wormTile, true),
			createPhaseTile(SmallTile, itileSpritesheet, wormTile, true),
		},
		[]PhaseEnemy{
			{
				EWorm,
				Vec2i{0, 0}, // Spawn on first tile
			},
			{
				EWorm,
				Vec2i{1, 1},
			},
			{
				EWorm,
				Vec2i{2, 0},
			},
			{
				EWorm,
				Vec2i{3, 1},
			},
		},
		[]string{
			"  ",
			"  ",
			"  ",
			"  ",
		},
	)
}
