package main

import (
	"image"

	paths "github.com/SolarLune/paths"
	"github.com/hajimehoshi/ebiten"
	pathfinding "github.com/xarg/gopathfinding"
)

// TODO: Make it so animations store if they play forwards or backwards in themselves
// TODO: Make it so that the beefy guy explodes before dying completely

// BeefEyeAnimations is the animations
type BeefEyeAnimations struct {
	idleFront Animation
	die       Animation
}

// BeefEyeAnimationSpeeds is the animation speeds
type BeefEyeAnimationSpeeds struct {
	idle float64
	die  float64
}

// BeefEye is a beefy eye type enemy
type BeefEye struct {
	position  Vec2f
	center    Vec2f
	size      Vec2i
	velocity  Vec2f
	moveSpeed float64

	health    int
	maxHealth int
	dead      bool
	remove    bool // Do we remove this enemy?
	flipped   bool // Is the enemy flipped?

	shadow *Shadow // The shadow below the enemy

	subImageRect image.Rectangle

	spritesheet     Spritesheet
	animation       Animation
	animations      BeefEyeAnimations
	animationSpeeds BeefEyeAnimationSpeeds

	image *ebiten.Image

	astarNodes  []pathfinding.Node
	path        paths.Path
	canPathfind bool
	pathFinding bool
}

func createBeefEye(position Vec2f, game *Game) *BeefEye {
	idleFrontSpritesheet := createSpritesheet(newVec2i(0, 23), newVec2i(234, 47), 9, ienemiesSpritesheet)
	dieSpritesheet := createSpritesheet(newVec2i(0, 48), newVec2i(570, 71), 19, ienemiesSpritesheet)

	shadowRect := image.Rect(0, 231, 14, 237)
	shadow := createShadow(shadowRect, iplayerSpritesheet, generateUniqueShadowID(game))
	game.shadows = append(game.shadows, &shadow)

	b := &BeefEye{
		position:  position,
		velocity:  newVec2f(0, 0),
		moveSpeed: 1.2,

		health:    6,
		maxHealth: 6,
		dead:      false,

		shadow: &shadow,

		spritesheet: dieSpritesheet,
		animations: BeefEyeAnimations{
			idleFront: createAnimation(idleFrontSpritesheet, ienemiesSpritesheet),
			die:       createAnimation(dieSpritesheet, ienemiesSpritesheet),
		},
		animationSpeeds: BeefEyeAnimationSpeeds{
			idle: 0.9,
			die:  1.2,
		},

		astarNodes:  []pathfinding.Node{},
		canPathfind: true,
		pathFinding: false,

		image: ienemiesSpritesheet,
	}

	return b
}

func (b *BeefEye) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// FLIP DECIDER
	flip := newVec2f(-1, 1)
	if b.flipped {
		flip.x = 1
	}

	// ROTATE & FLIP
	op.GeoM.Translate(float64(0-b.size.x)/2, float64(0-b.size.y)/2)
	op.GeoM.Scale(flip.x, flip.y)
	op.GeoM.Translate(float64(b.size.x)/2, float64(b.size.y)/2)
	b.subImageRect = image.Rect(
		b.spritesheet.sprites[b.animation.currentFrame].startPosition.x,
		b.spritesheet.sprites[b.animation.currentFrame].startPosition.y,
		b.spritesheet.sprites[b.animation.currentFrame].endPosition.x,
		b.spritesheet.sprites[b.animation.currentFrame].endPosition.y,
	)
	// POSITION
	op.GeoM.Translate(float64(b.position.x), float64(b.position.y))
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?

	screen.DrawImage(b.image.SubImage(b.subImageRect).(*ebiten.Image), op) // Draw enemy
}

func (b *BeefEye) update(game *Game) {
	// Start the animation if it's not playing
	b.animation = b.animations.idleFront
	b.animation.update(b.animationSpeeds.die)

	// Pathfind to player
	b.followPlayer(game)

	// Move enemy
	b.position.x += b.velocity.x
	b.position.y += b.velocity.y

	b.size = newVec2i(
		b.spritesheet.sprites[b.animation.currentFrame].size.x,
		b.spritesheet.sprites[b.animation.currentFrame].size.y,
	)
	endPosition := newVec2i(
		int(b.position.x)+b.size.x,
		int(b.position.y)+b.size.y,
	)

	// Update shadow
	b.shadow.update(b.position, b.size)

	b.subImageRect = image.Rect(int(b.position.x), int(b.position.y), endPosition.x, endPosition.y)
	b.center = newVec2f(b.position.x+float64(b.size.x)/2, b.position.y+float64(b.size.y)/2)
}

func (b *BeefEye) followPlayer(game *Game) {

	if b.canPathfind {
		start := newRolumn(
			int(b.position.x),
			int(b.position.y),
		)
		end := newRolumn(
			int(game.player.position.x),
			int(game.player.position.y),
		)

		// Make a path concurrently
		wg.Add(1)
		go calculatePath(astarChannel, game.currentMap.mapNodes, start, end)

		wg.Wait()

		// Get the path if it's finished
		b.path = *<-astarChannel
		b.canPathfind = false
	} else if !b.pathFinding && !b.canPathfind {

		if &b.path != nil {

			// If path is finished, generate new one.
			if len(b.path.Cells)-1 == b.path.CurrentIndex {
				b.canPathfind = true
			} else {
				z := b.size.x
				if b.size.y > b.size.x {
					z = b.size.y
				}
				ease := z // How much give the engine gives to 'reaching' a node

				finished := newVec2b(false, false)
				if len(b.path.Cells) > 0 {
					if int(b.position.x) < (b.path.Current().X*smallTileSize.x - ease) {
						b.position.x += b.moveSpeed
					} else if int(b.position.x) > (b.path.Current().X*smallTileSize.x + ease) {
						b.position.x -= b.moveSpeed
					} else {
						finished.x = true
					}

					if int(b.position.y) < (b.path.Current().Y*smallTileSize.y - ease) {
						b.position.y += b.moveSpeed
					} else if int(b.position.y) > (b.path.Current().Y*smallTileSize.y + ease) {
						b.position.y -= b.moveSpeed
					} else {
						finished.y = true
					}
				}

				if finished.x && finished.y {
					if b.path.AtEnd() {
						if b.path.AtEnd() {
							b.canPathfind = true
						}
					} else {
						if !b.path.AtEnd() {
							b.path.Next()
						}
					}
				}
			}
		}
	}
}

func (b *BeefEye) isDead() bool {
	if b.health <= 0 {
		if !b.dead {
			b.dead = true
		}
		return true
	}
	return false
}

func (b *BeefEye) attack(game *Game) {

}

func (b *BeefEye) getCenter() Vec2f {
	return b.center
}

func (b *BeefEye) getCurrentSubImageRect() image.Rectangle {
	return b.subImageRect
}

func (b *BeefEye) getImage() *ebiten.Image {
	return b.image
}

func (b *BeefEye) damage(value int) {
	b.health -= value
}

func (b *BeefEye) getShadow() Shadow {
	return *b.shadow
}
