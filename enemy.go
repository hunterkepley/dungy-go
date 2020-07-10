package main

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

// Enemy is the interface for all enemies in the game
type Enemy interface {
	render(screen *ebiten.Image)
	update()
	isDead() bool
	getCenter() Vec2f
	getCurrentSubImageRect() image.Rectangle
	getImage() *ebiten.Image
	damage()
}

func updateEnemies(g *Game) {
	// Temporarily spawn a worm
	if ebiten.IsKeyPressed(ebiten.KeyM) {
		g.enemies = append(g.enemies, Enemy(createWorm(newVec2f(float64(rand.Intn(screenWidth)), float64(rand.Intn(screenHeight))))))
	}

	for e := 0; e < len(g.enemies); e++ {
		if e >= len(g.enemies) {
			break
		}

		// Bullet collisions
		for i := 0; i < len(g.player.gun.bullets); i++ {

			if isAABBCollision(g.player.gun.bullets[i].collisionRect, g.enemies[e].getCurrentSubImageRect()) {
				g.enemies[e].damage()
				g.player.gun.bullets[i].destroy = true
			}
		}

		if g.enemies[e].isDead() {
			gibHandler := createGibHandler()
			gibHandler.explode(
				5,
				7,
				g.enemies[e].getCenter(),
				g.enemies[e].getCurrentSubImageRect(),
				g.enemies[e].getImage(),
			)
			g.gibHandlers = append(g.gibHandlers, gibHandler)
			g.enemies = removeEnemy(g.enemies, e)
			continue
		}
		g.enemies[e].update()

	}
}

func renderEnemies(g *Game, screen *ebiten.Image) {
	for _, e := range g.enemies {
		e.render(screen)
	}
}

func removeEnemy(slice []Enemy, e int) []Enemy {
	return append(slice[:e], slice[e+1:]...)
}
