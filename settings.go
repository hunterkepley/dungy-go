package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// SETTINGS FILE

// Settings is the game settings
type Settings struct {
	Graphics struct {
		Gibs       int  `yaml:"gibs"`
		Fullscreen bool `yaml:"fullscreen"`
	}
	Audio struct {
		MusicVolume float64 `yaml:"musicvolume:"`
	}
}

func loadSettings(s *Settings) {

	f, err := os.Open("./Assets/Config/settings.yaml")
	if err != nil {
		log.Printf("error -- Failed to load settings.yaml   #%v ", err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&s)
	if err != nil {
		log.Printf("error -- Failed to decode settings.yaml   #%v ", err)
	}

}
