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
	walkSide Animation
	die      Animation
}

// BeefEyeAnimationSpeeds is the animation speeds
type BeefEyeAnimationSpeeds struct {
	walk float64
	die  float64
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
	walkSideSpritesheet := createSpritesheet(newVec2i(0, 23), newVec2i(234, 47), 9, ienemiesSpritesheet)
	dieSpritesheet := createSpritesheet(newVec2i(0, 48), newVec2i(570, 71), 19, ienemiesSpritesheet)

	shadowRect := image.Rect(0, 231, 14, 237)
	shadow := createShadow(shadowRect, iplayerSpritesheet, generateUniqueShadowID(game))
	game.shadows = append(game.shadows, &shadow)

	b := &BeefEye{
		position:  position,
		velocity:  newVec2f(0, 0),
		moveSpeed: 1.2,

		health:                 15,
		maxHealth:              15,
		dead:                   false,
		idle:                   true,
		deathExplosion:         false,
		deathExplosionFinished: false,
		dying:                  false,

		shadow: &shadow,

		spritesheet: dieSpritesheet,
		animations: BeefEyeAnimations{
			walkSide: createAnimation(walkSideSpritesheet, ienemiesSpritesheet),
			die:      createAnimation(dieSpritesheet, ienemiesSpritesheet),
		},
		animationSpeeds: BeefEyeAnimationSpeeds{
			walk: 1.3,
			die:  2.,
		},

		astarNodes:  []pathfinding.Node{},
		canPathfind: true,
		pathFinding: false,

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
	}
	if b.deathExplosion && !b.dying {
		b.dying = true
		b.animation = b.animations.die
		b.animation.startBackwards()
	}

	switch b.animation.id {
	case b.animations.walkSide.id:
		b.animation.update(b.animationSpeeds.walk)
	case b.animations.die.id: // Special case for explosion animation
		b.animation.update(b.animationSpeeds.die)
		if b.animation.finishedFirstPlay {
			b.deathExplosionFinished = true
			b.deathExplosion = false
		}
	}

	if !b.dying {

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
