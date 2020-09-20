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
	dx := a.Min.X + a.Size().X/2 - b.Min.X + b.Size().X/2
	dy := a.Min.Y + a.Size().Y/2 - b.Min.Y + b.Size().Y/2
	distance := math.Sqrt(float64(dx*dx + dy*dy))

	r1 := a.Size().X
	if a.Size().Y > r1 {
		r1 = a.Size().Y
	}
	r2 := b.Size().X
	if b.Size().Y > r2 {
		r2 = b.Size().Y
	}
	if distance < float64(r1+r2) {
		return true
	}
	return false
}
