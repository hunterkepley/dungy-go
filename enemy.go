package main

import (
	"image"

	paths "github.com/SolarLune/paths"
	"github.com/hajimehoshi/ebiten"
)

// EnemyType is the type of an enemy (Beefeye, Worm, etc)
type EnemyType int

const (
	// EBeefEye is the Beef Eye enemy
	EBeefEye EnemyType = iota + 1

	// EWorm is the worm enemy
	EWorm
)

// Enemy is the interface for all enemies in the game
type Enemy interface {
	render(screen *ebiten.Image)
	update(game *Game)
	isDead() bool
	damage(value int)
	attack(game *Game)

	setPosition(Vec2f)

	getShadow() Shadow
	getCenter() Vec2f
	getPosition() Vec2f
	getCurrentSubImageRect() image.Rectangle // Subimage rect
	getImage() *ebiten.Image                 // Image
	getFlipped() bool                        // Is enemy flipped?
	getSize() Vec2i
	getMoveSpeed() float64
	getDying() bool     // Is the enemy dying?
	getAttacking() bool // Is the enemy attacking currently?

	// Mostly pathfinding stuff
	getCanPathfind() bool // Can enemy find a path?
	getPath() *paths.Path // The enemy's path
	getPathfinding() bool // Is the enemy pathfinding currently?

	setFlipped(bool)
	setPath(paths.Path)
	setCanPathfind(bool)
}

// Generates an enemy then adds to enemy list automatically
func generateEnemy(enemyType EnemyType, position Vec2f, g *Game) {
	switch enemyType {
	case EBeefEye:
		g.enemies = append(g.enemies, createBeefEye(position, g))
	case EWorm:
		g.enemies = append(g.enemies, createWorm(position, g))
	}
}

func updateEnemies(g *Game) {

	for e := 0; e < len(g.enemies); e++ {
		if e >= len(g.enemies) {
			break
		}

		if g.enemies[e].isDead() {
			gibAmount := 8 // Gib setting 1
			gibSize := 7   // Gib setting 1

			switch g.settings.Graphics.Gibs { // Gib setting 2
			case 2:
				gibSize = 7
				gibAmount = 15
			case 0:
				gibAmount = 0
			}

			gibHandler := createGibHandler()
			g.shadows = removeShadow((g.shadows), g.enemies[e].getShadow().id)

			gibHandler.explode(
				gibAmount,
				gibSize,
				g.enemies[e].getCenter(),
				g.enemies[e].getCurrentSubImageRect(),
				g.enemies[e].getImage(),
			)

			g.gibHandlers = append(g.gibHandlers, gibHandler)
			g.enemies = removeEnemy(g.enemies, e)
			continue
		}
		g.enemies[e].update(g)

		if !g.enemies[e].getDying() && !g.enemies[e].getAttacking() {
			enemiesPathfinding(g, g.enemies[e])
		}

		// Bullet collisions
		for b := 0; b < len(g.player.gun.bullets); b++ {
			if isAABBCollision(g.enemies[e].getCurrentSubImageRect(), g.player.gun.bullets[b].collisionRect) {
				g.enemies[e].damage(g.player.gun.calculateDamage())
				// TODO: This causes a very annoying bug where light appears at top right for a frame sometimes
				g.player.gun.bullets[b].destroy = true
				break
			}
		}

	}
}

func enemiesPathfinding(g *Game, e Enemy) {
	if g.player.position.x+float64(g.player.dynamicSize.x)/2 > e.getCenter().x {
		e.setFlipped(true)
	} else {
		e.setFlipped(false)
	}

	if e.getCanPathfind() {
		start := newRolumn(
			int(e.getPosition().x),
			int(e.getPosition().y),
		)
		end := newRolumn(
			int(g.player.position.x),
			int(g.player.position.y),
		)

		// Make a path concurrently
		wg.Add(1)
		go calculatePath(astarChannel, g.currentMap.mapNodes, start, end)

		defer wg.Wait()

		// Get the path if it's finished
		e.setPath(*<-astarChannel)
		e.setCanPathfind(false)
	} else if !e.getPathfinding() && !e.getCanPathfind() {

		if e.getPath() != nil {

			// If path is finished, generate new one.
			if len(e.getPath().Cells)-1 == e.getPath().CurrentIndex {
				e.setCanPathfind(true)
			} else {
				z := e.getSize().x
				if e.getSize().y > e.getSize().x {
					z = e.getSize().y
				}
				ease := z // How much give the engine gives to 'reaching' a node

				finished := newVec2b(false, false)
				if len(e.getPath().Cells) > 0 {
					if int(e.getPosition().x) < (e.getPath().Current().X*smallTileSize.x - ease) {
						e.setPosition(Vec2f{e.getPosition().x + e.getMoveSpeed(), e.getPosition().y})
					} else if int(e.getPosition().x) > (e.getPath().Current().X*smallTileSize.x + ease) {
						e.setPosition(Vec2f{e.getPosition().x - e.getMoveSpeed(), e.getPosition().y})
					} else {
						finished.x = true
					}

					if int(e.getPosition().y) < (e.getPath().Current().Y*smallTileSize.y - ease) {
						e.setPosition(Vec2f{e.getPosition().x, e.getPosition().y + e.getMoveSpeed()})
					} else if int(e.getPosition().y) > (e.getPath().Current().Y*smallTileSize.y + ease) {
						e.setPosition(Vec2f{e.getPosition().x, e.getPosition().y - e.getMoveSpeed()})
					} else {
						finished.y = true
					}
				}

				if finished.x && finished.y {
					if e.getPath().AtEnd() {
						if e.getPath().AtEnd() {
							e.setCanPathfind(true)
						}
					} else {
						if !e.getPath().AtEnd() {
							e.getPath().Next()
						}
					}
				}
			}
		}
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

func enemySpawnHandler(g *Game) {

}
