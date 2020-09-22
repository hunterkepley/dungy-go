package main

import (
	paths "github.com/SolarLune/paths"
)

func calculatePath(game *Game, id int, mapNodes []string, start Rolumn, end Rolumn) {
	// This line creates a new Grid, comprised of Cells. The size is 10x10. By default, all Cells are
	// walkable and have a cost of 1, and a blank character of ' '.
	//firstMap := paths.NewGrid(10, 10)

	mapLayout := paths.NewGridFromStringArrays(mapNodes)

	// After creating the Grid, you can edit it using the Grid's functions. Note that here, we're using 'x'
	// to get Cells that have the rune for the lowercase x character 'x', not the string "x".
	for _, cell := range mapLayout.GetCellsByRune('x') {
		cell.Walkable = false
	}

	for _, goop := range mapLayout.GetCellsByRune('g') {
		goop.Cost = 5
	}

	// This gets a new Path (a slice of Cells) from the starting Cell to the destination Cell. If the path's length
	// is greater than 0, then it was successful.
	path := mapLayout.GetPath(mapLayout.Get(start.row, start.column), mapLayout.Get(end.row, end.column), true)

	game.astarChannels[id] <- path
}
