package main

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

// Player is the player in the fuckin game I hate useless comments
type Player struct {
	position Vec2f
	image    *ebiten.Image
	speed    float64
}

func createPlayer(position Vec2f) Player {
	speed := 5.
	return Player{
		position,
		playerSpritesheet,
		speed,
	}
}

func (p *Player) update() {
	p.input()
}

func (p *Player) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.position.x, p.position.y)
	//i := (g.count / 5) % frameNum
	//sx, sy := frameOX+i*frameWidth, frameOY
	//image, _ := loadPicture("./Assets/Art/Player/rightIdle.png")
	screen.DrawImage(playerSpritesheet.SubImage(image.Rect(0, 0, 15, 15)).(*ebiten.Image), op) //runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
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
