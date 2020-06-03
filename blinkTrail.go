package main

// BlinkTrail the trail when the player blinks
type BlinkTrail struct {
	speed    float64
	speedMax float64

	sections []BlinkTrailSection
}

// BlinkTrailSection is the actual animated section of the trail
type BlinkTrailSection struct {
	position Vec2f

	spritesheet Spritesheet
	animation   Animation
}
