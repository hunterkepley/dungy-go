package main

import "github.com/hajimehoshi/ebiten"

// Bullet is for the player's gun to fire
type Bullet struct {
	position Vec2f
	size     Vec2i

	image *ebiten.Image
}

func (b *Bullet) render(screen *ebiten.Image) {

}

func (b *Bullet) update() {

}
