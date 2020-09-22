package main

import (
	"fmt"

	paths "github.com/SolarLune/paths"
)

func calculatePath(id int, mapNodes []string, start Rolumn, end Rolumn) {
	// This line creates a new Grid, comprised of Cells. The size is 10x10. By default, all Cells are
	// walkable and have a cost of 1, and a blank character of ' '.
	//firstMap := paths.NewGrid(10, 10)

	start = newRolumn(start.column/(len(mapNodes[0])-1), start.row/(len(mapNodes)-2))
	end = newRolumn(end.column/(len(mapNodes[0])-1), end.row/(len(mapNodes)-2))
	fmt.Println("s: ", start.column, ", ", start.row)
	fmt.Println("e: ", end.column, ", ", end.row)

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

	fmt.Println(len(astarChannels[id]))
	astarChannels[id] <- *mapLayout.GetPath(mapLayout.Get(start.column, start.row), mapLayout.Get(end.column, end.row), true)
}
