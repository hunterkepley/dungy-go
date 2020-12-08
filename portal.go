package main

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten"
)

// Portal is a portal to another world
type Portal struct {
	position Vec2f

	state int // 0 -> sucking tiles; 1 -> spewing tiles

	mapNumber int

	sprite Sprite
	image  *ebiten.Image
}

func createPortal(position Vec2f, mapNumber int) Portal {
	image := itileSpritesheet
	return Portal{
		position: position,

		state: -1,

		mapNumber: mapNumber,

		sprite: createSprite(Vec2i{21, 302}, Vec2i{34, 328}, Vec2i{13, 26}, image),
		image:  image,
	}
}

func (p *Portal) update(g *Game) {
	if isAABBCollision(
		image.Rect(int(p.position.x), int(p.position.y), int(p.position.x)+p.sprite.size.x, int(p.position.y)+p.sprite.size.y),
		g.player.getBoundsStatic()) && p.state == -1 {

		p.state = 0
	}
	switch p.state {
	case 0:
		p.eatTiles(g)
		g.player.position = p.position
	case 1:

	}
}

func (p *Portal) eatTiles(g *Game) {
	moveSpeed := 1.
	for i := 0; i < len(g.currentMap.tiles); i++ {
		for j := 0; j < len(g.currentMap.tiles[i]); j++ {
			// Calculate movement using an imaginary vector :)
			dx := p.position.x - g.currentMap.tiles[i][j].position.x
			dy := p.position.y - g.currentMap.tiles[i][j].position.y

			ln := math.Sqrt(dx*dx + dy*dy)

			dx /= ln
			dy /= ln

			// Move towards portal
			g.currentMap.tiles[i][j].position.x += dx * moveSpeed
			g.currentMap.tiles[i][j].position.y += dy * moveSpeed

			g.currentMap.tiles[i][j].rotation = dx * dy // Rotation
			if g.currentMap.tiles[i][j].scale.x > 0 {
				g.currentMap.tiles[i][j].scale.x -= 0.004 // Scale
			}
			if g.currentMap.tiles[i][j].scale.y > 0 {
				g.currentMap.tiles[i][j].scale.y -= 0.004 // Scale
			}

			if g.currentMap.tiles[i][j].scale.x <= 0 || g.currentMap.tiles[i][j].scale.y <= 0 {
				// Switch map
				g.currentMap = g.maps[1]
				p.state = 1
			}
		}
	}
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
