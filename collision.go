package main

import (
	"image"
	"math"
)

// isAABBCollision checks whether or not a collides with b using rectangles
func isAABBCollision(a image.Rectangle, b image.Rectangle) bool {
	if a.Min.X < b.Max.X &&
		a.Max.X > b.Min.X &&
		a.Min.Y < b.Max.Y &&
		a.Max.Y > b.Min.Y {

		return true
	}
	return false
}

// isCircularCollision checks whether or not a collides with b using circles based on the radius of the largest w/h of the rects
func isCircularCollision(a image.Rectangle, b image.Rectangle) bool {

	c1 := Vec2i{a.Min.X + (a.Max.X-a.Min.X)/2, a.Min.Y + (a.Max.Y-a.Min.Y)/2}
	c2 := Vec2i{b.Min.X + (b.Max.X-b.Min.X)/2, b.Min.Y + (b.Max.Y-b.Min.Y)/2}

	distance := math.Sqrt(float64(((c1.x - c2.x) * (c1.x - c2.x)) + ((c1.y - c2.y) * (c1.y - c2.y))))

	r1 := a.Max.X - a.Min.X
	if a.Max.Y-a.Min.Y > r1 {
		r1 = a.Max.Y - a.Min.Y
	}

	r2 := b.Max.X - b.Min.X
	if b.Max.Y-b.Min.Y > r1 {
		r2 = b.Max.Y - b.Min.Y
	}
	r2 /= 2
	r1 /= 2

	if distance < float64(r1+r2) {
		return true
	}

	return false
}
