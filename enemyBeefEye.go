package main

import (
	"image"

	paths "github.com/SolarLune/paths"
	"github.com/hajimehoshi/ebiten"
	pathfinding "github.com/xarg/gopathfinding"
)

// BeefEyeShockwave is each individual shockwave that can hit the player
type BeefEyeShockwave struct {
	position Vec2f

	velocity  Vec2f
	moveSpeed Vec2f

	animation          Animation
	animationSpeed     float64
	animationCanFinish bool

	image *ebiten.Image
}

// BeefEyeShockwaveHandler controls the shockwave that is emitted when the BeefEye attacks
type BeefEyeShockwaveHandler struct {
	shockwaves []BeefEyeShockwave
}

// BeefEyeAnimations is the animations
type BeefEyeAnimations struct {
	walkSide Animation
	die      Animation
	attack   Animation
}

func (b *BeefEyeShockwaveHandler) createShockwave(beefEye *BeefEye) {
	numberOfWaves := 10

	for i := 0; i < numberOfWaves; i++ {
		b.shockwaves = append(b.shockwaves, BeefEyeShockwave{
			position: beefEye.center,

			moveSpeed: newVec2f(10, 10),

			image: ienemiesSpritesheet,
		})
	}
}

// BeefEyeAnimationSpeeds is the animation speeds
type BeefEyeAnimationSpeeds struct {
	walk   float64
	die    float64
	attack float64
}

// BeefEye is a beefy eye type enemy
type BeefEye struct {
	position  Vec2f
	center    Vec2f
	size      Vec2i
	velocity  Vec2f
	moveSpeed float64

	health                 int
	maxHealth              int
	dead                   bool
	deathExplosion         bool // When the death explosion is playing
	deathExplosionFinished bool // When the death explosion is finished
	dying                  bool
	remove                 bool // Do we remove this enemy?
	flipped                bool // Is the enemy flipped?
	idle                   bool // Is the enemy idling?
	attacking              bool // Is the enemy attacking?

	shadow           *Shadow                 // The shadow below the enemy
	shockwaveHandler BeefEyeShockwaveHandler // The shockwave handler for when the BeefEye attacks

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
	walkSideSpritesheet := createSpritesheet(newVec2i(0, 23), newVec2i(234, 47), 9, ienemiesSpritesheet)
	dieSpritesheet := createSpritesheet(newVec2i(0, 48), newVec2i(570, 71), 19, ienemiesSpritesheet)
	attackSpritesheet := createSpritesheet(newVec2i(270, 48), newVec2i(570, 71), 10, ienemiesSpritesheet)

	shadowRect := image.Rect(0, 231, 14, 237)
	shadow := createShadow(shadowRect, iplayerSpritesheet, generateUniqueShadowID(game))
	game.shadows = append(game.shadows, &shadow)

	b := &BeefEye{
		position:  position,
		velocity:  newVec2f(0, 0),
		moveSpeed: 1.1,

		health:    15,
		maxHealth: 15,
		idle:      true,

		shadow: &shadow,

		spritesheet: dieSpritesheet,
		animations: BeefEyeAnimations{
			walkSide: createAnimation(walkSideSpritesheet, ienemiesSpritesheet),
			die:      createAnimation(dieSpritesheet, ienemiesSpritesheet),
			attack:   createAnimation(attackSpritesheet, ienemiesSpritesheet),
		},
		animationSpeeds: BeefEyeAnimationSpeeds{
			walk:   1.3,
			die:    2.,
			attack: 1.,
		},

		astarNodes:  []pathfinding.Node{},
		canPathfind: true,

		image: ienemiesSpritesheet,
	}

	return b
}

func (b *BeefEye) render(screen *ebiten.Image) {
	if len(b.animation.spritesheet.sprites) == 0 {
		return
	}
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
		b.animation.spritesheet.sprites[b.animation.currentFrame].startPosition.x,
		b.animation.spritesheet.sprites[b.animation.currentFrame].startPosition.y,
		b.animation.spritesheet.sprites[b.animation.currentFrame].endPosition.x,
		b.animation.spritesheet.sprites[b.animation.currentFrame].endPosition.y,
	)
	// POSITION
	op.GeoM.Translate(float64(b.position.x), float64(b.position.y))
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?

	screen.DrawImage(b.image.SubImage(b.subImageRect).(*ebiten.Image), op) // Draw enemy
}

func (b *BeefEye) update(game *Game) {
	// Start the animation if it's not playing
	if b.idle {
		b.animation = b.animations.walkSide
		b.animation.startForwards()
		b.idle = false
	} else if b.deathExplosion && !b.dying {
		b.dying = true
		b.animation = b.animations.die
		b.animation.startBackwards()
	}

	// Animation speeds and other specialities related to animations
	switch b.animation.id {
	case b.animations.walkSide.id:
		b.animation.update(b.animationSpeeds.walk)
	case b.animations.die.id: // Special case for explosion animation
		b.animation.update(b.animationSpeeds.die)
		if b.animation.finishedFirstPlay {
			b.deathExplosionFinished = true
			b.deathExplosion = false
		}
	case b.animations.attack.id:
		b.animation.update(b.animationSpeeds.attack)
	}

	if !b.dying && !b.attacking {

		// Move enemy
		b.position.x += b.velocity.x
		b.position.y += b.velocity.y

		b.size = newVec2i(
			b.animation.spritesheet.sprites[b.animation.currentFrame].size.x,
			b.animation.spritesheet.sprites[b.animation.currentFrame].size.y,
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

	// Attack

	b.attack(game)
}

func (b *BeefEye) isDead() bool {
	if b.health <= 0 {

		if !b.deathExplosionFinished {
			b.deathExplosion = true
		} else {
			if !b.dead && !b.deathExplosion {
				b.dead = true
			}
			return true
		}
	}
	return false
}

func (b *BeefEye) attack(game *Game) {
	if !b.dying {
		if b.attacking {
			if b.animation.finishedFirstPlay {
				b.attacking = false
				b.idle = true
			}
		} else if !b.attacking && isAABBCollision(game.player.getBoundsDynamic(), b.subImageRect) {

			b.animation = b.animations.attack
			b.animation.startBackwards()

			b.attacking = true

		}
	}
}

func (b *BeefEye) getCenter() Vec2f {
	return b.center
}

func (b *BeefEye) getPosition() Vec2f {
	return b.position
}

func (b *BeefEye) getCurrentSubImageRect() image.Rectangle {
	return b.subImageRect
}

func (b *BeefEye) getImage() *ebiten.Image {
	return b.image
}

func (b *BeefEye) getSize() Vec2i {
	return b.size
}

func (b *BeefEye) damage(value int) {
	b.health--
}

func (b *BeefEye) getShadow() Shadow {
	return *b.shadow
}

func (b *BeefEye) getMoveSpeed() float64 {
	return b.moveSpeed
}

func (b *BeefEye) getDying() bool {
	return b.dying
}

func (b *BeefEye) getAttacking() bool {
	return b.attacking
}

func (b *BeefEye) getPath() *paths.Path {
	return &b.path
}

func (b *BeefEye) getCanPathfind() bool {
	return b.canPathfind
}

func (b *BeefEye) getFlipped() bool {
	return b.flipped
}

func (b *BeefEye) getPathfinding() bool {
	return b.pathFinding
}

func (b *BeefEye) setPosition(pos Vec2f) {
	b.position = pos
}

func (b *BeefEye) setFlipped(flipped bool) {
	b.flipped = flipped
}

func (b *BeefEye) setPath(path paths.Path) {
	b.path = path
}

func (b *BeefEye) setCanPathfind(canPathfind bool) {
	b.canPathfind = canPathfind
}
