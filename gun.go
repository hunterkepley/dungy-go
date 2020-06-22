package main

import (
	"math"

	"github.com/hajimehoshi/ebiten"
)

// Gun is the player's gun
type Gun struct {
	position Vec2f
	size     Vec2i

	rotation float64

	image *ebiten.Image
}

func (g *Gun) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	op.GeoM.Rotate(g.rotation)
	op.GeoM.Translate(g.position.x, g.position.y)
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?
	screen.DrawImage(g.image, op)
}

func (g *Gun) update(playerPosition Vec2f) {
	g.position = playerPosition
	/*mouseX, mouseY = pygame.mouse.get_pos()
	 *playerX, playerY = player.get_pos()
	 *angle = math.atan2(playerX-mouseX, playerY-mouseY)
	 */
	cursorPosition := Vec2i{}
	cursorPosition.x, cursorPosition.y = ebiten.CursorPosition()
	g.rotation = math.Atan2(g.position.x-float64(cursorPosition.x), g.position.y-float64(cursorPosition.y))
}
