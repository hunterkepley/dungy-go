package main

// MapData is the metadata for a Map (mainly the name and version)
type MapData struct {
	name    string
	version string
}

// Map is a map, it contains MapData, tiles, etc
type Map struct {
	data MapData

	tiles []Tile
}
