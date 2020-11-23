package main

import (
	"github.com/hajimehoshi/ebiten"
)

// Portal is a portal to another world
type Portal struct {
	position Vec2f

	sprite Sprite
	image  *ebiten.Image
}

func createPortal(position Vec2f) Portal {
	image := itileSpritesheet
	return Portal{
		position: position,

		sprite: createSprite(Vec2i{21, 302}, Vec2i{34, 328}, Vec2i{13, 26}, image),
		image:  image,
	}
}

func (p *Portal) update(g *Game) {

}

func (p *Portal) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.position.x, p.position.y)
	op.Filter = ebiten.FilterNearest

	// TODO: Animate
	//currentFrame := w.animation.spritesheet.sprites[w.animation.currentFrame]

	/*subImageRect := image.Rect(
		currentFrame.startPosition.x,
		currentFrame.startPosition.y,
		currentFrame.endPosition.x,
		currentFrame.endPosition.y,
	)*/

	screen.DrawImage(p.image.SubImage(p.sprite.getBounds()).(*ebiten.Image), op) // Draw player
}
