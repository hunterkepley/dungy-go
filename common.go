package main

// Vec2f is a vector of 2 float64's
type Vec2f struct {
	x float64
	y float64
}

// NewVec2f creates a new Vec2f and returns it
func NewVec2f(x float64, y float64) Vec2f {
	return Vec2f{x, y}
}
