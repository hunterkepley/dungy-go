package main

import "fmt"

import "github.com/nickdavies/go-astar/astar"

func main () {
    rows := 3
    cols := 3

    // Build AStar object from existing
    // PointToPoint configuration
    a := astar.NewAStar(rows, cols)
    p2p := astar.NewPointToPoint()

    // Make an invincible obsticle at (1,1)
    a.FillTile(astar.Point{1, 1}, -1) 

    // Path from one corner to the other
    source := []astar.Point{astar.Point{0,0}}
    target := []astar.Point{astar.Point{2,2}}

    path := a.FindPath(p2p, source, target)

    for path != nil {
        fmt.Printf("At (%d, %d)\n", path.Col, path.Row)
        path = path.Parent
    }
}