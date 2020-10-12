package main

import "math/rand"

// PhaseChunk is a chunk of a "map" to phase in randomly
type PhaseChunk struct {
	tiles    []Tile
	enemies  []Enemy
	mapNodes []string
}

func createPhaseChunk(tiles []Tile, enemies []Enemy, mapNodes []string) PhaseChunk {
	return PhaseChunk{
		tiles:    tiles,
		enemies:  enemies,
		mapNodes: mapNodes,
	}
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
		p.timer = p.timerMax
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

		timer:    500,
		timerMax: 500,
	}

	p.makePhases()

	return p
}

func (p *Phases) makePhases() {
	p.addPhaseChunk(
		[]Tile{},
		[]Enemy{},
		[]string{},
	)
}

func (p *Phases) getRandomPhase() PhaseChunk {
	return p.chunks[rand.Intn(len(p.chunks))]
}

func (p *Phases) addPhaseChunk(tiles []Tile, enemies []Enemy, mapNodes []string) {
	c := createPhaseChunk(tiles, enemies, mapNodes)
	p.chunks = append(p.chunks, c)
}
