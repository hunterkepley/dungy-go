package main

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

// Player is the player in the fuckin game I hate useless comments
type Player struct {
	position Vec2f
	//center   Vec2f
	speed       float64
	spritesheet Spritesheet // Current spritesheet
	animation   Animation   // Current animation
	image       *ebiten.Image
}

func createPlayer(position Vec2f) Player {
	speed := 3.
	image := playerSpritesheet
	spritesheet := createSpritesheet(newVec2i(0, 0), newVec2i(75, 26), 5, image)
	return Player{
		position,
		speed,
		spritesheet,
		createAnimation(spritesheet, image),
		image,
	}
}

func (p *Player) update() {
	p.animation.play(0.3)
	p.input()
}

func (p *Player) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.position.x, p.position.y)
	subImageRect := image.Rect(
		p.spritesheet.sprites[p.animation.currentFrame].startPosition.x,
		p.spritesheet.sprites[p.animation.currentFrame].startPosition.y,
		p.spritesheet.sprites[p.animation.currentFrame].endPosition.x,
		p.spritesheet.sprites[p.animation.currentFrame].endPosition.y,
	)
	screen.DrawImage(p.image.SubImage(subImageRect).(*ebiten.Image), op)
}

func (p *Player) input() {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.position.y -= p.speed
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.position.y += p.speed
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.position.x -= p.speed
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.position.x += p.speed
	}
}
