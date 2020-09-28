package main

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v2"
)

// SETTINGS FILE

// Settings is the game settings
type Settings struct {
	Graphics struct {
		gibs       int  // 0 -> off, 1 -> normal, 2 -> high
		fullscreen bool // true for on, false for off
	}
}

var testSettings = `
graphics:
	gibs: 2
	fullscreen: false
`

func (s *Settings) loadSettings() Settings {
	var settings Settings

	err := yaml.Unmarshal([]byte(testSettings), &settings)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- m:\n%v\n\n", settings)
	return settings
}
