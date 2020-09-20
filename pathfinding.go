package main

import (
	"fmt"

	"github.com/nickdavies/go-astar/astar"
	pathfinding "github.com/xarg/gopathfinding"
)

// NodeType is a type for the node struct
type NodeType int

const (
	// Walkable ... NODETYPE ENUM [1]
	Walkable NodeType = iota + 1
	// Unwalkable ... NODETYPE ENUM [2]
	Unwalkable
)

func (n NodeType) String() string {
	return [...]string{"Unknown", "Walkable", "Unwalkable"}[n]
}

// Node is a astar node
type Node struct {
	rolumn Rolumn
	_type  NodeType
}

// AStarContext holds the different parts of the astar system for scope
type AStarContext struct {
	a        astar.AStar
	p2p      astar.AStarConfig
	pathChan *chan astar.PathPoint
}

func createAStarContext(rows int, columns int) AStarContext {
	// Build AStar object from existing
	// PointToPoint configuration
	a := astar.NewAStar(rows, columns)
	p2p := astar.NewPointToPoint()

	return AStarContext{
		a:   a,
		p2p: p2p,
	}
}

func initAStar(game *Game, context *AStarContext, nodes []Node) {
	// Unwalkable nodes
	for i := 0; i < len(nodes); i++ {
		switch nodes[i]._type {
		case Unwalkable:
			x := nodes[i].rolumn.column
			y := nodes[i].rolumn.row
			context.a.FillTile(astar.Point{Row: x, Col: y}, -1)
		case Walkable:
			// nothin to do here for now
		}
	}
}

func calculatePath(game *Game, id int, nodes []Node, start Rolumn, end Rolumn) astar.PathPoint {
	p := astar.PathPoint{}
	//A pathfinding.MapData containing the
	//coordinates(x, y) of LAND, WALL, START and STOP of the map.
	//If your map is something more than 2d matrix then you might want to modify adjacentNodes

	graph := pathfinding.NewGraph(map_data)

	//Returns a list of nodes from START to STOP avoiding all obstacles if possible
	shortest_path := pathfinding.Astar(graph)
	if &game.astarContexts[id] != nil {

		/*go calculatePoints(
		context,
		, ,
		)*/

		a := astarContext

		fmt.Println(a)

		s := []astar.Point{{Row: start.row, Col: start.column}}
		t := []astar.Point{{Row: end.row, Col: end.column}}
		p = *a.a.FindPath(a.p2p, s, t)

		//*context.pathChan <- *p

	}
	return p
}
