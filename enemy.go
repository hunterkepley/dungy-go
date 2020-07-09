package main

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

// Enemy is the interface for all enemies in the game
type Enemy interface {
	render(screen *ebiten.Image)
	update(bullets []Bullet)
	isDead() bool
	getCenter() Vec2f
	getCurrentSubImageRect() image.Rectangle
	getImage() *ebiten.Image
}

func updateEnemies(g *Game) {
	// Temporarily spawn a worm
	if ebiten.IsKeyPressed(ebiten.KeyM) {
		g.enemies = append(g.enemies, Enemy(createWorm(newVec2f(float64(rand.Intn(screenWidth)), float64(rand.Intn(screenHeight))))))
	}
	//fmt.Println(len(g.enemies))
	for i, e := range g.enemies {
		if i >= len(g.enemies) {
			break
		}

		if e.isDead() {
			gibHandler := createGibHandler()
			gibHandler.explode(5, 7, e.getCenter(), e.getCurrentSubImageRect(), e.getImage())
			g.gibHandlers = append(g.gibHandlers, gibHandler)
			g.enemies = remove(g.enemies, i)
			continue
		}
		e.update(g.player.gun.bullets)

	}
}

func renderEnemies(g *Game, screen *ebiten.Image) {
	for _, e := range g.enemies {
		e.render(screen)
	}
}

func remove(slice []Enemy, e int) []Enemy {
	return append(slice[:e], slice[e+1:]...)
}
