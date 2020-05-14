package main

// Vec2f is a vector of 2 float64's
type Vec2f struct {
	x float64
	y float64
}

// NewVec2f creates a new Vec2f and returns it
func newVec2f(x float64, y float64) Vec2f {
	return Vec2f{x, y}
}

// Vec2i is a vector 2 2 int's
type Vec2i struct {
	x int
	y int
}

// NewVec2i creates a new Vec2i and returns it
func newVec2i(x int, y int) Vec2i {
	return Vec2i{x, y}
}
