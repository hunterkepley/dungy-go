package main

import (
	"sync"

	paths "github.com/SolarLune/paths"
)

var (
	wg sync.WaitGroup

	astarChannel chan *paths.Path // Stores all channels for astar paths
)

/*

Possible solution for waitgroup based crashes:

Make a counter for how many concurrent processes are finished or running, then only
reuse the waitgroup when when Wait has returned?

Or multiple waitgroups? Every enemy could have one

*/

func calculatePath(channel chan *paths.Path, mapNodes []string, start Rolumn, end Rolumn, enemy *Enemy) {
	// This line creates a new Grid, comprised of Cells. The size is 10x10. By default, all Cells are
	// walkable and have a cost of 1, and a blank character of ' '.
	//firstMap := paths.NewGrid(10, 10)

	defer wg.Done()

	start = newRolumn(start.column/smallTileSize.x, start.row/smallTileSize.y)
	end = newRolumn(end.column/smallTileSize.x, end.row/smallTileSize.y)

	mapLayout := paths.NewGridFromStringArrays(mapNodes)

	// Take enemies into account for pathfinding
	for _, e := range gameReference.enemies {
		pos := Vec2i{int(e.getPosition().x) / smallTileSize.x, int(e.getPosition().y) / smallTileSize.y}

		if e.getShadow().id != (*enemy).getShadow().id {

			if len(mapLayout.Data) > pos.y && len(mapLayout.Data[0]) > pos.x {
				mapLayout.Data[pos.y][pos.x] = &paths.Cell{X: pos.x, Y: pos.y, Character: 'x'}
			}
		}

	}

	// After creating the Grid, you can edit it using the Grid's functions. Note that here, we're using 'x'
	// to get Cells that have the rune for the lowercase x character 'x', not the string "x".
	for _, cell := range mapLayout.GetCellsByRune('x') {
		cell.Walkable = false
	}

	for _, goop := range mapLayout.GetCellsByRune('g') {
		goop.Cost = 3
	}

	for _, lava := range mapLayout.GetCellsByRune('l') {
		lava.Cost = 15
	}

	// This gets a new Path (a slice of Cells) from the starting Cell to the destination Cell. If the path's length
	// is greater than 0, then it was successful.
	for {
		// Break out with a nil path if unable to path (start/end on non-walkable cells)
		if !mapLayout.Get(start.column, start.row).Walkable || !mapLayout.Get(end.column, end.row).Walkable {
			channel <- mapLayout.GetPath(mapLayout.Get(3, 3), mapLayout.Get(3, 4), false)
			return
		}
		p := mapLayout.GetPath(mapLayout.Get(start.column, start.row), mapLayout.Get(end.column, end.row), false)
		if p != nil {
			channel <- p
			break
		}
	}
}
