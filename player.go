package main

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

// Player is the player in the fuckin game I hate useless comments
type Player struct {
	position Vec2f
	image    *ebiten.Image
}

func createPlayer(position Vec2f) Player {
	return Player{
		position,
		playerSpritesheet,
	}
}

func (p *Player) update() {

}

func (p *Player) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	//i := (g.count / 5) % frameNum
	//sx, sy := frameOX+i*frameWidth, frameOY
	//image, _ := loadPicture("./Assets/Art/Player/rightIdle.png")
	screen.DrawImage(playerSpritesheet.SubImage(image.Rect(0, 0, 15, 15)).(*ebiten.Image), op) //runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
}
