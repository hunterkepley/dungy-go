package main

import (
	"fmt"
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

	astarChannelID         int
	astarNodes             []pathfinding.Node
	pathChan               *chan *paths.Path
	path                   *paths.Path
	pathfindingTickRate    int
	pathfindingTickRateMax int
	pathfindingIndex       paths.Cell
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

		astarNodes:             []pathfinding.Node{},
		pathfindingTickRate:    5,
		pathfindingTickRateMax: 5,

		image: ienemiesSpritesheet,
	}

	// Pathfinding stuff
	b.astarChannelID = len(game.astarChannels) - 1

	game.astarChannels = append(game.astarChannels, make(chan *paths.Path, 2000))
	b.pathChan = &game.astarChannels[b.astarChannelID]

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

	if b.path != nil && b.path.AtEnd() {
		numTiles := newVec2i(getNumberOfTilesPossible().x, getNumberOfTilesPossible().y)
		start := newRolumn(
			int(b.position.x)/(numTiles.x-2),
			int(b.position.y)/(numTiles.y-2),
		)
		end := newRolumn(
			int(game.player.position.x)/(numTiles.x-2),
			int(game.player.position.y)/(numTiles.y-2),
		)

		go calculatePath(game, b.astarChannelID, game.currentMap.mapNodes, start, end)

		select {
		case b.path = <-*b.pathChan:
		default:
		}
	} else {

		if &b.pathfindingIndex == nil && b.path != nil {
			b.pathfindingIndex = *b.path.Current()
		}

		if &b.pathfindingIndex != nil {

			ease := 10 // How much give the engine gives to 'reaching' a node

			finished := newVec2b(false, false)
			if int(b.position.x) < b.pathfindingIndex.X-ease {
				b.position.x += b.moveSpeed
			} else if int(b.position.x) > b.pathfindingIndex.X+ease {
				b.position.x -= b.moveSpeed
			} else {
				finished.x = true
			}

			if int(b.position.y) < b.pathfindingIndex.Y-ease {
				b.position.y += b.moveSpeed
			} else if int(b.position.y) > b.pathfindingIndex.Y+ease {
				b.position.y -= b.moveSpeed
			} else {
				finished.y = true
			}

			if finished.x && finished.y {
				fmt.Println("Node complete at: ", b.position.x, ", ", b.position.y)
				if b.path.Next() == nil {
					b.path.Next()

				} else {
					b.pathfindingIndex = *b.path.Next()
				}
			} else {
				fmt.Println("Node unfinished")
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
