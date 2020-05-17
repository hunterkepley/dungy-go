package main

// Direction is a type for the direction enum [for player/enemies]
type Direction int

const (
	// Up ... DIRECTION ENUM [1]
	Up Direction = iota + 1
	// Down ... DIRECTION ENUM [2]
	Down
	// Left ... DIRECTION ENUM [3]
	Left
	// Right ... DIRECTION ENUM [4]
	Right
)

func (d Direction) String() string {
	return [...]string{"Unknown", "Up", "Down", "Left", "Right"}[d]
}

// ^ DIRECTION ENUM ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Movement is a type for the movement enum [for player/enemies]
type Movement int

const (
	// Idle ... MOVEMENT ENUM [1]
	Idle Movement = iota + 1
	// Walking ... MOVEMENT ENUM [2]
	Walking
	// Running ... MOVEMENT ENUM [3]
	Running
)

func (m Movement) String() string {
	return [...]string{"Unknown", "Idle", "Walking", "Running"}[m]
}

// ^ MOVEMENT ENUM ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

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
