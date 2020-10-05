package main

import "math/rand"

// PhaseChunk is a chunk of a "map" to phase in randomly
type PhaseChunk struct {
	tiles    []Tile
	mapNodes []string
}

/*
 * This function handles the phase changes, or, the chunks of new maps that replace the current one
 *
 * Lots to do
 */
func (m *Map) phaseHandler() {

}

// Phases is a struct containing all the phases in the game and some controls for them
type Phases struct {
	chunks []PhaseChunk
}

func (p *Phases) getRandomPhase() PhaseChunk {
	return p.chunks[rand.Intn(len(p.chunks))]
}

func (p *Phases) addPhaseChunk() {
	c := PhaseChunk{}
	p.chunks = append(p.chunks, c)
}
