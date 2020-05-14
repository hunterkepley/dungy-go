package main

// Direction is a type for the direction enum [for player/enemies]
type Direction int

const (
	// Up ... DIRECTION ENUM [0]
	Up Direction = iota
	// Down ... DIRECTION ENUM [1]
	Down
	// Left ... DIRECTION ENUM [2]
	Left
	// Right ... DIRECTION ENUM [3]
	Right
)

func (d Direction) String() string {
	return [...]string{"Up", "Down", "Left", "Right"}[d]
}

// ^ DIRECTION ENUM ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Vec2f is a vector of 2 float64's
type Vec2f struct {
	x float64
	y float64
}

func newVec2f(x float64, y float64) Vec2f {
	return Vec2f{x, y}
}

// Vec2i is a vector 2 2 int's
type Vec2i struct {
	x int
	y int
}

func newVec2i(x int, y int) Vec2i {
	return Vec2i{x, y}
}
