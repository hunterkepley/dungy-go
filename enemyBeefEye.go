package main

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten"
	"github.com/nickdavies/go-astar/astar"
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

	findPath       bool
	astarContextID int
	astarNodes     []Node
	pathChan       *chan astar.PathPoint
	path           astar.PathPoint
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
		moveSpeed: 0.4,

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

		findPath:   false,
		astarNodes: []Node{},

		image: ienemiesSpritesheet,
	}

	// Pathfinding stuff
	astarContext := createAStarContext(getNumberOfTilesPossible().x, getNumberOfTilesPossible().y)
	game.astarContexts = append(game.astarContexts, astarContext)
	b.astarContextID = len(game.astarContexts) - 1

	initAStar(game, &game.astarContexts[b.astarContextID], b.astarNodes)

	game.astarChannels = append(game.astarChannels, make(chan astar.PathPoint, 2000))
	b.pathChan = &game.astarChannels[b.astarContextID]
	game.astarContexts[b.astarContextID].pathChan = b.pathChan

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
	numTiles := newVec2i(getNumberOfTilesPossible().x, getNumberOfTilesPossible().y)
	if !b.findPath {

		start := newRolumn(
			int(b.position.x)/numTiles.x,
			int(b.position.y)/numTiles.y,
		)
		end := newRolumn(
			int(game.player.position.x)/numTiles.x,
			int(game.player.position.y)/numTiles.y,
		)

		b.path = calculatePath(game, b.astarContextID, b.astarNodes, start, end)

		/*if b.pathChan != nil {
			select {
			case b.path = <-*b.pathChan: // Path made
				fmt.Println("hello")
				b.findPath = true
				if &b.path != nil {
					fmt.Println("Calculated")
				} else {
					fmt.Println("FUCK")
					b.findPath = false
				}
			default: // No path made
			}
		}*/
		if &b.path != nil {
			b.findPath = false
		}
	} else {
		iterator := b.path
		for &b.path != nil {
			b.position = newVec2f(float64(iterator.Col*numTiles.x), float64(iterator.Row*numTiles.y))
			fmt.Println(b.position.x, ", ", b.position.y)
			if b.path.Parent != nil {
				iterator = *b.path.Parent
			} else {
				break
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
